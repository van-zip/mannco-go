package mannco

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/cart/get" {
			t.Errorf("Expected /cart/get, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		response := APIResponse[CartPayload]{
			Err:     false,
			Success: true,
			Content: CartPayload{
				Cart: []CartItem{
					{
						CartID:  42,
						AssetID: "987654321,987654322",
						Count:   2,
						ItemID:  5678,
						Price:   15000,
						Name:    "Unusual Burning Flames Team Captain",
						Game:    440,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-jwt", nil)
	client.SetBaseURL(server.URL)

	ctx := context.Background()
	cart, err := client.GetCart(ctx)
	if err != nil {
		t.Fatalf("GetCart failed: %v", err)
	}

	if len(cart.Cart) != 1 {
		t.Fatalf("Expected 1 cart item, got %d", len(cart.Cart))
	}

	item := cart.Cart[0]
	if item.CartID != 42 {
		t.Errorf("Expected cartId 42, got %d", item.CartID)
	}
	if item.AssetID != "987654321,987654322" {
		t.Errorf("Expected assetId '987654321,987654322', got %s", item.AssetID)
	}
	if item.Count != 2 {
		t.Errorf("Expected count 2, got %d", item.Count)
	}
	if item.ItemID != 5678 {
		t.Errorf("Expected item_id 5678, got %d", item.ItemID)
	}
	if item.Price != 15000 {
		t.Errorf("Expected price 15000, got %d", item.Price)
	}
	if item.Name != "Unusual Burning Flames Team Captain" {
		t.Errorf("Expected name 'Unusual Burning Flames Team Captain', got %s", item.Name)
	}
	if item.Game != 440 {
		t.Errorf("Expected game 440, got %d", item.Game)
	}
}

func TestAddToCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cart/add" {
			t.Errorf("Expected /cart/add, got %s", r.URL.Path)
		}

		var req AddToCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		if req.AssetID != "987654321" {
			t.Errorf("Expected assetId '987654321', got %s", req.AssetID)
		}

		w.Header().Set("Content-Type", "application/json")
		response := APIResponse[CartPayload]{
			Err:     false,
			Success: true,
			Content: CartPayload{
				Cart: []CartItem{
					{
						CartID:  43,
						AssetID: "987654321",
						Count:   1,
						ItemID:  5678,
						Price:   15000,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-jwt", nil)
	client.SetBaseURL(server.URL)

	ctx := context.Background()
	cart, err := client.AddToCart(ctx, "987654321")
	if err != nil {
		t.Fatalf("AddToCart failed: %v", err)
	}

	if len(cart.Cart) != 1 {
		t.Fatalf("Expected 1 cart item, got %d", len(cart.Cart))
	}

	item := cart.Cart[0]
	if item.CartID != 43 {
		t.Errorf("Expected cartId 43, got %d", item.CartID)
	}
	if item.AssetID != "987654321" {
		t.Errorf("Expected assetId '987654321', got %s", item.AssetID)
	}
}

func TestBulkAddToCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cart/bulk" {
			t.Errorf("Expected /cart/bulk, got %s", r.URL.Path)
		}

		var req BulkAddToCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		if req.ItemID != 5678 {
			t.Errorf("Expected itemId 5678, got %d", req.ItemID)
		}
		if req.Count != 5 {
			t.Errorf("Expected count 5, got %d", req.Count)
		}
		if req.SellerUserID != "76561198000000000" {
			t.Errorf("Expected sellerUserId '76561198000000000', got %s", req.SellerUserID)
		}

		w.Header().Set("Content-Type", "application/json")
		response := APIResponse[CartPayload]{
			Err:     false,
			Success: true,
			Content: CartPayload{
				Cart: []CartItem{
					{
						CartID:  44,
						AssetID: "111,222,333",
						Count:   3,
						ItemID:  5678,
						Price:   14000,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-jwt", nil)
	client.SetBaseURL(server.URL)

	ctx := context.Background()
	cart, err := client.BulkAddToCart(ctx, 5678, 5, "76561198000000000")
	if err != nil {
		t.Fatalf("BulkAddToCart failed: %v", err)
	}

	if len(cart.Cart) != 1 {
		t.Fatalf("Expected 1 cart item, got %d", len(cart.Cart))
	}

	item := cart.Cart[0]
	if item.CartID != 44 {
		t.Errorf("Expected cartId 44, got %d", item.CartID)
	}
	if item.AssetID != "111,222,333" {
		t.Errorf("Expected assetId '111,222,333', got %s", item.AssetID)
	}
	if item.Count != 3 {
		t.Errorf("Expected count 3, got %d", item.Count)
	}
	if item.Price != 14000 {
		t.Errorf("Expected price 14000, got %d", item.Price)
	}
}

func TestRemoveFromCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cart/remove" {
			t.Errorf("Expected /cart/remove, got %s", r.URL.Path)
		}

		var req RemoveFromCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		if req.CartID != 42 {
			t.Errorf("Expected cartId 42, got %d", req.CartID)
		}

		w.Header().Set("Content-Type", "application/json")
		response := APIResponse[CartPayload]{
			Err:     false,
			Success: true,
			Content: CartPayload{
				Cart: []CartItem{},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-jwt", nil)
	client.SetBaseURL(server.URL)

	ctx := context.Background()
	cart, err := client.RemoveFromCart(ctx, 42)
	if err != nil {
		t.Fatalf("RemoveFromCart failed: %v", err)
	}

	if len(cart.Cart) != 0 {
		t.Fatalf("Expected empty cart, got %d items", len(cart.Cart))
	}
}

func TestUpdateCart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cart/update" {
			t.Errorf("Expected /cart/update, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		response := APIResponse[CartUpdateResponse]{
			Err:     false,
			Success: true,
			Content: CartUpdateResponse{
				Integrity: struct {
					Valid        bool               `json:"valid"`
					InvalidItems []CartIntegrityItem `json:"invalidItems"`
				}{
					Valid: false,
					InvalidItems: []CartIntegrityItem{
						{
							CartID:        42,
							AssetID:       "987654321",
							Reason:        "price_changed",
							CurrentPrice:  16000,
							ExpectedPrice: 15000,
						},
					},
				},
				Replaced: []CartReplacedItem{
					{
						CartID:     42,
						OldAssetID: "987654321",
						NewAssetID: "987654999",
						Reason:     "price_changed",
					},
				},
				Removed: []int{},
				Cart: []CartItem{
					{
						CartID:  42,
						AssetID: "987654999",
						Count:   1,
						ItemID:  5678,
						Price:   15000,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-jwt", nil)
	client.SetBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.UpdateCart(ctx)
	if err != nil {
		t.Fatalf("UpdateCart failed: %v", err)
	}

	if result.Integrity.Valid {
		t.Error("Expected integrity valid to be false")
	}

	if len(result.Integrity.InvalidItems) != 1 {
		t.Fatalf("Expected 1 invalid item, got %d", len(result.Integrity.InvalidItems))
	}

	invalid := result.Integrity.InvalidItems[0]
	if invalid.CartID != 42 {
		t.Errorf("Expected cartId 42, got %d", invalid.CartID)
	}
	if invalid.Reason != "price_changed" {
		t.Errorf("Expected reason 'price_changed', got %s", invalid.Reason)
	}
	if invalid.CurrentPrice != 16000 {
		t.Errorf("Expected currentPrice 16000, got %d", invalid.CurrentPrice)
	}
	if invalid.ExpectedPrice != 15000 {
		t.Errorf("Expected expectedPrice 15000, got %d", invalid.ExpectedPrice)
	}

	if len(result.Replaced) != 1 {
		t.Fatalf("Expected 1 replaced item, got %d", len(result.Replaced))
	}

	replaced := result.Replaced[0]
	if replaced.CartID != 42 {
		t.Errorf("Expected cartId 42, got %d", replaced.CartID)
	}
	if replaced.OldAssetID != "987654321" {
		t.Errorf("Expected oldAssetId '987654321', got %s", replaced.OldAssetID)
	}
	if replaced.NewAssetID != "987654999" {
		t.Errorf("Expected newAssetId '987654999', got %s", replaced.NewAssetID)
	}

	if len(result.Cart) != 1 {
		t.Fatalf("Expected 1 cart item, got %d", len(result.Cart))
	}

	item := result.Cart[0]
	if item.CartID != 42 {
		t.Errorf("Expected cartId 42, got %d", item.CartID)
	}
	if item.AssetID != "987654999" {
		t.Errorf("Expected assetId '987654999', got %s", item.AssetID)
	}
	if item.Price != 15000 {
		t.Errorf("Expected price 15000, got %d", item.Price)
	}
}