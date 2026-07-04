package mannco

import (
	"context"
	"encoding/json"
	"fmt"
)

// InventoryItem represents a single item in the user's inventory or sale listings.
type InventoryItem struct {
	IDs    string `json:"ids"`
	Count  int    `json:"count"`
	ItemID int    `json:"item_id"`
	State  int    `json:"state"`
	Price  int    `json:"price,omitempty"`
	Name   string `json:"name"`
	Game   int    `json:"game"`
}

// InventoryItemsPayload is the response payload for inventory listing endpoints.
type InventoryItemsPayload struct {
	Items []InventoryItem `json:"items"`
	Count int             `json:"count"`
}

// WithdrawItemsResponse is the response from a withdrawal request.
type WithdrawItemsResponse struct {
	Message string `json:"message"`
	Updated int    `json:"updated"`
	Locked  int    `json:"locked"`
}

// GetItemsOnSale returns all items currently listed for sale by the authenticated user.
func (c *Client) GetItemsOnSale(ctx context.Context) (InventoryItemsPayload, error) {
	return executeRequest[InventoryItemsPayload](ctx, c, "GET", "inventory/onSale", nil, nil)
}

// GetItemsInInventory returns all items in the authenticated user's inventory not listed for sale.
func (c *Client) GetItemsInInventory(ctx context.Context) (InventoryItemsPayload, error) {
	return executeRequest[InventoryItemsPayload](ctx, c, "GET", "inventory/onInventory", nil, nil)
}

// SetItemPrice sets a price for items and matches them with existing buy orders.
func (c *Client) SetItemPrice(ctx context.Context, ids string, price int) error {
	payload := struct {
		IDs   string `json:"ids"`
		Price int    `json:"price"`
	}{
		IDs:   ids,
		Price: price,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%w: error encoding json for set item price: %w", ErrInternal, err)
	}
	_, err = executeRequest[json.RawMessage](ctx, c, "POST", "inventory/price", jsonData, nil)
	return err
}

// WithdrawItems initiates a trade offer to withdraw items from the user's inventory.
func (c *Client) WithdrawItems(ctx context.Context, ids string) (WithdrawItemsResponse, error) {
	payload := struct {
		IDs string `json:"ids"`
	}{
		IDs: ids,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return WithdrawItemsResponse{}, fmt.Errorf("%w: error encoding json for withdraw items: %w", ErrInternal, err)
	}
	return executeRequest[WithdrawItemsResponse](ctx, c, "POST", "inventory/withdraw", jsonData, nil)
}