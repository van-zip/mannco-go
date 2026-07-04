package mannco

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// UserLogin uses the provided API key to login and return a JWT
func (c *Client) UserLogin(ctx context.Context, apiKey string) (string, error) {
	data := map[string]string{"apiKey": apiKey}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error encoding json for user login: %w", err)
	}

	content, err := executeRequest[LoginPayload](ctx, c, "POST", "user/login", jsonData, nil)
	if err != nil {
		return "", err
	}

	c.SetJWT(content.JWT)

	return content.JWT, nil
}

// Balance returns the user balance in pennies
func (c *Client) Balance(ctx context.Context) (int, error) {
	content, err := executeRequest[BalancePayload](ctx, c, "GET", "user/balance", nil, nil)
	if err != nil {
		return 0, err
	}
	return content.Balance, nil
}

// TransactionHistory returns the user transaction history on the site
func (c *Client) TransactionHistory(ctx context.Context, opts *HistoryOptions) (InventoryPayload, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
	}
	return executeRequest[InventoryPayload](ctx, c, "GET", "user/getTransactionHistory", nil, params)
}

// SalesHistory returns the user sales history on the site
func (c *Client) SalesHistory(ctx context.Context, opts *HistoryOptions) (InventoryPayload, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("perPage", strconv.Itoa(opts.Limit))
		}
		if opts.Period != "" {
			params.Add("timeRange", string(opts.Period))
		}
		if opts.Search != "" {
			params.Add("search", opts.Search)
		}
	}
	return executeRequest[InventoryPayload](ctx, c, "GET", "user/getSalesHistory", nil, params)
}

// PurchaseHistory returns the user purchase history on the site
func (c *Client) PurchaseHistory(ctx context.Context, opts *HistoryOptions) (InventoryPayload, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			params.Add("count", strconv.Itoa(opts.Limit))
		}
	}
	return executeRequest[InventoryPayload](ctx, c, "GET", "user/getPurchaseHistory", nil, params)
}

// UserInfo contains the authenticated user's account details.
type UserInfo struct {
	SteamID      string `json:"steamId"`
	Balance      int    `json:"balance"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	TradeURL     string `json:"tradeurl"`
	TwoFA        string `json:"2fa"`
	ShortURL     string `json:"shorturl"`
	Notification int    `json:"notification"`
}

// UserInfoPayload wraps the user info response.
type UserInfoPayload struct {
	Informations UserInfo `json:"informations"`
}

// NotificationPayload is the response from /user/notifications.
type NotificationPayload struct {
	AlertCount    int `json:"alertCount"`
	MessagesCount int `json:"messagesCount"`
	OffersCount   int `json:"offersCount"`
	TradesCount   int `json:"tradesCount"`
}

// BalanceHistoryPayload is the response from /user/getBalanceHistory.
type BalanceHistoryPayload struct {
	Values []json.RawMessage `json:"values"`
	Count  int               `json:"count"`
}

// CashoutHistoryPayload is the response from /user/getCashoutHistory.
type CashoutHistoryPayload struct {
	Values []json.RawMessage `json:"values"`
}

// TransactionDetailsPayload is the response from /user/getTransactionDetails.
type TransactionDetailsPayload struct {
	Transaction json.RawMessage `json:"transaction"`
}

// StoreProfile is the response from /user/store/{identifier}.
type StoreProfile struct {
	SteamID string `json:"steamId"`
	Name    string `json:"name"`
	Image   string `json:"image"`
}

// SessionPayload is the response from /user/ipList.
type SessionPayload struct {
	Values []json.RawMessage `json:"values"`
	Count  int               `json:"count"`
}

// DisconnectPayload is the response from /user/disconnect.
type DisconnectPayload struct {
	Disconnect bool `json:"disconnect"`
}

// UserInfo returns account details for the authenticated user.
func (c *Client) UserInfo(ctx context.Context) (UserInfoPayload, error) {
	return executeRequest[UserInfoPayload](ctx, c, "GET", "user/infos", nil, nil)
}

// Notifications returns aggregated unread notification counts.
func (c *Client) Notifications(ctx context.Context) (NotificationPayload, error) {
	return executeRequest[NotificationPayload](ctx, c, "GET", "user/notifications", nil, nil)
}

// BalanceHistory returns paginated balance change records.
func (c *Client) BalanceHistory(ctx context.Context, page, limit int) (BalanceHistoryPayload, error) {
	params := url.Values{}
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	return executeRequest[BalanceHistoryPayload](ctx, c, "GET", "user/getBalanceHistory", nil, params)
}

// CashoutHistory returns paginated cashout records.
func (c *Client) CashoutHistory(ctx context.Context, page, limit int) (CashoutHistoryPayload, error) {
	params := url.Values{}
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	return executeRequest[CashoutHistoryPayload](ctx, c, "GET", "user/getCashoutHistory", nil, params)
}

// TransactionDetails returns the detail payload for a single transaction.
func (c *Client) TransactionDetails(ctx context.Context, transactionID string) (TransactionDetailsPayload, error) {
	params := url.Values{}
	params.Add("transactionId", transactionID)
	return executeRequest[TransactionDetailsPayload](ctx, c, "GET", "user/getTransactionDetails", nil, params)
}

// StoreProfile returns public store profile information by store identifier.
func (c *Client) StoreProfile(ctx context.Context, identifier string) (StoreProfile, error) {
	return executeRequest[StoreProfile](ctx, c, "GET", "user/store/"+identifier, nil, nil)
}

// SalesStats returns summary sales metrics for the authenticated user.
func (c *Client) SalesStats(ctx context.Context) (json.RawMessage, error) {
	return executeRequest[json.RawMessage](ctx, c, "GET", "user/getSalesInfos", nil, nil)
}

// SalesChart returns chart-oriented sales data.
func (c *Client) SalesChart(ctx context.Context, period, chart string) (json.RawMessage, error) {
	params := url.Values{}
	if period != "" {
		params.Add("period", period)
	}
	if chart != "" {
		params.Add("chart", chart)
	}
	return executeRequest[json.RawMessage](ctx, c, "GET", "user/getSalesChartInfos", nil, params)
}

// SessionList returns IP/session rows for the authenticated user.
func (c *Client) SessionList(ctx context.Context, page, perPage int, includeExpired bool) (SessionPayload, error) {
	params := url.Values{}
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Add("perPage", strconv.Itoa(perPage))
	}
	if includeExpired {
		params.Add("expire", "true")
	}
	return executeRequest[SessionPayload](ctx, c, "GET", "user/ipList", nil, params)
}

// Disconnect invalidates the current session / JWT token.
func (c *Client) Disconnect(ctx context.Context) error {
	content, err := executeRequest[DisconnectPayload](ctx, c, "GET", "user/disconnect", nil, nil)
	if err != nil {
		return err
	}
	if !content.Disconnect {
		return fmt.Errorf("%w: disconnect returned false", ErrInternal)
	}
	return nil
}
