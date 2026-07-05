// Package mannco is an API client for mannco.store
package mannco

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// BaseURL is the default base URL for the mannco.store api
const BaseURL = "https://api.mannco.store/"

// Retry configuration constants
const (
	MaxRetries       = 3
	BaseRetryDelay   = 1 * time.Second
	MaxRetryDelay    = 30 * time.Second
	DefaultRetryWait = 60 * time.Second
)

// APIResponse is the general shape of Mannco.store API responses
type APIResponse[T any] struct {
	Err     bool   `json:"err"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Content T      `json:"content"`
}

// APIError represents the general shape of API errors
type APIError struct {
	StatusCode int
	Message    string
}

// Error formats APIErrors
func (e *APIError) Error() string {
	return fmt.Sprintf("server error with status code %d: %s", e.StatusCode, e.Message)
}

var (
	// ErrUnauthorized indicates either insufficient permissions or expired authentication
	ErrUnauthorized = errors.New("authentication error")
	// ErrNetwork is any network related failure
	ErrNetwork = errors.New("network error")
	// ErrInternal is any error derived from internal package logic
	ErrInternal = errors.New("internal error")
)

// Client represents the base API client which interacts with Mannco.store
type Client struct {
	httpClient *http.Client
	mu         sync.RWMutex
	baseURL    string
	jwt        string
	apiKey     string
}

// NewClient instantiates a new API client
func NewClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 60 * time.Second,
		}
	}
	return &Client{
		baseURL:    BaseURL,
		jwt:        "", // until login this isn't set
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// SetJWT sets the JWT for a client
func (c *Client) SetJWT(token string) {
	c.mu.Lock()
	c.jwt = token
	c.mu.Unlock()
}

// GetJWT gets the JWT for a client
func (c *Client) GetJWT() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.jwt
}

// SetAPIKey sets the API key for a client
func (c *Client) SetAPIKey(key string) {
	c.mu.Lock()
	c.apiKey = key
	c.mu.Unlock()
}

// GetAPIKey gets the API key for a client
func (c *Client) GetAPIKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiKey
}

// SetBaseURL sets the API client base url
func (c *Client) SetBaseURL(url string) {
	c.mu.Lock()
	c.baseURL = url
	c.mu.Unlock()
}

// GetBaseURL gets the API client base url
func (c *Client) GetBaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL

}

// executeRequest performs generic parsing, safety handling, and raw IO operations for interacting with the API
func executeRequest[T any](ctx context.Context, c *Client, method, endpoint string, body []byte, queryParams url.Values) (T, error) {
	var target T

	u, err := url.Parse(c.GetBaseURL() + endpoint)
	if err != nil {
		return target, fmt.Errorf("%w: invalid endpoint url: %w", ErrInternal, err)
	}
	if queryParams != nil {
		u.RawQuery = queryParams.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return target, fmt.Errorf("%w: failed to construct request object: %w", ErrInternal, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if jwt := c.GetJWT(); jwt != "" {
		req.Header.Set("Authorization", "Bearer "+jwt)
	}

	// Execute with retry logic for 429 status
	// these are closed already so we can safely ignore it
	resp, bodyBytes, err := executeWithRetry(ctx, c, req) //nolint:bodyclose
	if err != nil {
		return target, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		if endpoint != "user/login" {
			err := c.UserLogin(ctx)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+c.GetJWT())
				resp, bodyBytes, err = executeWithRetry(ctx, c, req) //nolint:bodyclose
				if err != nil {
					return target, err
				}
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIResponse[json.RawMessage]
		msg := ""
		if json.Unmarshal(bodyBytes, &apiErr) == nil {
			msg = apiErr.Message
		}
		httpErr := &APIError{StatusCode: resp.StatusCode, Message: msg}
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return target, fmt.Errorf("%w: %w", ErrUnauthorized, httpErr)
		}
		return target, fmt.Errorf("%w: %w", ErrInternal, httpErr)
	}
	var apiResponse APIResponse[T]
	if err = json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return target, fmt.Errorf("%w: failed decoding response JSON: %w", ErrInternal, &APIError{StatusCode: resp.StatusCode, Message: err.Error()})
	}

	if apiResponse.Err || !apiResponse.Success {
		return target, fmt.Errorf("%w: %w", ErrInternal, &APIError{StatusCode: resp.StatusCode, Message: apiResponse.Message})
	}
	return apiResponse.Content, nil
}

// executeWithRetry executes an HTTP request with retry logic for 429 (rate limited) responses
func executeWithRetry(ctx context.Context, c *Client, req *http.Request) (*http.Response, []byte, error) {
	var resp *http.Response
	var bodyBytes []byte
	var err error

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: request execution failed: %w", ErrNetwork, err)
		}

		bodyBytes, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, nil, fmt.Errorf("%w: failed reading raw response bytes: %w", ErrNetwork, &APIError{StatusCode: resp.StatusCode, Message: err.Error()})
		}

		// success
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, bodyBytes, nil
		}

		if attempt == MaxRetries {
			return resp, bodyBytes, nil // Return the 429 response on final attempt
		}

		delay := time.Duration(attempt+1) * BaseRetryDelay
		if delay > MaxRetryDelay {
			delay = MaxRetryDelay
		}

		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(delay):
		}
	}

	return resp, bodyBytes, nil
}
