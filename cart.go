package mannco

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// CartItem represents a cart row grouping one or more asset IDs of the same item at the same price.
type CartItem struct {
	CartID  int    `json:"cartId"`
	AssetID string `json:"assetId"` // Comma-separated Steam asset IDs
	Count   int    `json:"count"`
	ItemID  int    `json:"item_id"`
	Price   int    `json:"price"` // Price in cents
	Name    string `json:"name"`
	Game    int    `json:"game"`
}

// CartPayload wraps the cart array response
type CartPayload struct {
	Cart []CartItem `json:"cart"`
}

// AssetIDs returns the individual Steam asset IDs split from the comma-separated field.
func (c CartItem) AssetIDs() []string {
	if c.AssetID == "" {
		return nil
	}
	return strings.Split(c.AssetID, ",")
}

// AddToCartRequest for POST /cart/add
type AddToCartRequest struct {
	AssetID string `json:"assetId"`
}

// BulkAddToCartRequest for POST /cart/bulk
type BulkAddToCartRequest struct {
	ItemID       int    `json:"itemId"`
	Count        int    `json:"count"`
	SellerUserID string `json:"sellerUserId,omitempty"`
}

// RemoveFromCartRequest for POST /cart/remove
type RemoveFromCartRequest struct {
	CartID int `json:"cartId"`
}

// CartIntegrityItem represents an invalid item found during integrity check
type CartIntegrityItem struct {
	CartID        int    `json:"cartId"`
	AssetID       string `json:"assetId"`
	Reason        string `json:"reason"`         // e.g., "price_changed"
	CurrentPrice  int    `json:"currentPrice"`
	ExpectedPrice int    `json:"expectedPrice"`
}

// CartReplacedItem represents an item that was replaced during update
type CartReplacedItem struct {
	CartID     int    `json:"cartId"`
	OldAssetID string `json:"oldAssetId"`
	NewAssetID string `json:"newAssetId"`
	Reason     string `json:"reason"`
}

// CartUpdateResponse for POST /cart/update
type CartUpdateResponse struct {
	Integrity struct {
		Valid        bool               `json:"valid"`
		InvalidItems []CartIntegrityItem `json:"invalidItems"`
	} `json:"integrity"`
	Replaced []CartReplacedItem `json:"replaced"`
	Removed  []int              `json:"removed"`
	Cart     []CartItem         `json:"cart"`
}

// GetCart retrieves the authenticated user's cart
func (c *Client) GetCart(ctx context.Context) (CartPayload, error) {
	return executeRequest[CartPayload](ctx, c, "GET", "cart/get", nil, nil)
}

// AddToCart adds a single item to the cart by asset ID
func (c *Client) AddToCart(ctx context.Context, assetID string) (CartPayload, error) {
	payload := AddToCartRequest{AssetID: assetID}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return CartPayload{}, fmt.Errorf("%w: error encoding json for add to cart: %w", ErrInternal, err)
	}
	return executeRequest[CartPayload](ctx, c, "POST", "cart/add", jsonData, nil)
}

// BulkAddToCart adds multiple items of the same type, selecting cheapest available listings
func (c *Client) BulkAddToCart(ctx context.Context, itemID, count int, sellerUserID string) (CartPayload, error) {
	payload := BulkAddToCartRequest{
		ItemID:       itemID,
		Count:        count,
		SellerUserID: sellerUserID,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return CartPayload{}, fmt.Errorf("%w: error encoding json for bulk add to cart: %w", ErrInternal, err)
	}
	return executeRequest[CartPayload](ctx, c, "POST", "cart/bulk", jsonData, nil)
}

// RemoveFromCart removes a cart row by cart ID
func (c *Client) RemoveFromCart(ctx context.Context, cartID int) (CartPayload, error) {
	payload := RemoveFromCartRequest{CartID: cartID}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return CartPayload{}, fmt.Errorf("%w: error encoding json for remove from cart: %w", ErrInternal, err)
	}
	return executeRequest[CartPayload](ctx, c, "POST", "cart/remove", jsonData, nil)
}

// UpdateCart runs integrity checks and optionally fixes the cart
func (c *Client) UpdateCart(ctx context.Context) (CartUpdateResponse, error) {
	return executeRequest[CartUpdateResponse](ctx, c, "POST", "cart/update", []byte("{}"), nil)
}
