package mannco

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestCreatePayment(t *testing.T) {
	runAPITest(t, testCase[PaymentResponse]{
		name:           "CreatePayment_balance_redirect",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"url":"https://payment-provider.example/checkout/session_abc"}}`,
		expectedPath:   "/payment/payviox",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (PaymentResponse, error) {
			return client.CreatePayment(ctx, ProviderPayviox, CreatePaymentRequest{
				Type:         PaymentTypeBalance,
				TOSTimestamp: 1706745600,
				Value:        10000,
			})
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req CreatePaymentRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.Type != PaymentTypeBalance {
				t.Errorf("Expected type balance, got %s", req.Type)
			}
			if req.TOSTimestamp != 1706745600 {
				t.Errorf("Expected tos_timestamp 1706745600, got %d", req.TOSTimestamp)
			}
			if req.Value != 10000 {
				t.Errorf("Expected value 10000, got %d", req.Value)
			}
			if req.Items != "" {
				t.Errorf("Expected empty items, got %s", req.Items)
			}
		},
		assertResponse: func(t *testing.T, resp PaymentResponse) {
			if resp.URL != "https://payment-provider.example/checkout/session_abc" {
				t.Errorf("Expected redirect URL, got %s", resp.URL)
			}
			if resp.Message != "" {
				t.Errorf("Expected no message for redirect, got %s", resp.Message)
			}
		},
	})

	runAPITest(t, testCase[PaymentResponse]{
		name:           "CreatePayment_items_processed",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"message":"Payment processed successfully"}}`,
		expectedPath:   "/payment/mannco",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (PaymentResponse, error) {
			return client.CreatePayment(ctx, ProviderMannco, CreatePaymentRequest{
				Type:         PaymentTypeItems,
				TOSTimestamp: 1706745600,
				Items:        "987654321,987654322",
			})
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req CreatePaymentRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.Type != PaymentTypeItems {
				t.Errorf("Expected type items, got %s", req.Type)
			}
			if req.TOSTimestamp != 1706745600 {
				t.Errorf("Expected tos_timestamp 1706745600, got %d", req.TOSTimestamp)
			}
			if req.Items != "987654321,987654322" {
				t.Errorf("Expected items '987654321,987654322', got %s", req.Items)
			}
			if req.Value != 0 {
				t.Errorf("Expected value 0 for items payment, got %d", req.Value)
			}
		},
		assertResponse: func(t *testing.T, resp PaymentResponse) {
			if resp.Message != "Payment processed successfully" {
				t.Errorf("Expected success message, got %s", resp.Message)
			}
			if resp.URL != "" {
				t.Errorf("Expected no URL for processed payment, got %s", resp.URL)
			}
		},
	})

	runAPITest(t, testCase[PaymentResponse]{
		name:           "CreatePayment_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"content":"forbidden"}`,
		expectedPath:   "/payment/payviox",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (PaymentResponse, error) {
			return client.CreatePayment(ctx, ProviderPayviox, CreatePaymentRequest{
				Type:         PaymentTypeBalance,
				TOSTimestamp: 1706745600,
				Value:        10000,
			})
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})

	runAPITest(t, testCase[PaymentResponse]{
		name:           "CreatePayment_business_error",
		mockStatus:     300,
		mockResponse:   `{"err":true,"success":false,"content":"Insufficient balance"}`,
		expectedPath:   "/payment/mannco",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (PaymentResponse, error) {
			return client.CreatePayment(ctx, ProviderMannco, CreatePaymentRequest{
				Type:         PaymentTypeItems,
				TOSTimestamp: 1706745600,
				Items:        "987654321",
			})
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrInternal) {
				t.Errorf("Expected ErrInternal, got %v", err)
			}
		},
	})
}

