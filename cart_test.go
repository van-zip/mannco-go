package mannco

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestGetCart(t *testing.T) {
	runAPITest(t, testCase[CartPayload]{
		name:           "GetCart_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"cart":[{"cartId":42,"assetId":"987654321,987654322","count":2,"item_id":5678,"price":15000,"name":"Unusual Burning Flames Team Captain","game":440}]}}`,
		expectedPath:   "/cart/get",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.GetCart(ctx)
		},
		assertResponse: func(t *testing.T, cart CartPayload) {
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
		},
	})

	runAPITest(t, testCase[CartPayload]{
		name:           "GetCart_unauthorized",
		mockStatus:     401,
		mockResponse:   `{"err":true,"success":false,"message":"Unauthorized","content":null}`,
		expectedPath:   "/cart/get",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.GetCart(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}

func TestAddToCart(t *testing.T) {
	runAPITest(t, testCase[CartPayload]{
		name:           "AddToCart_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"cart":[{"cartId":43,"assetId":"987654321","count":1,"item_id":5678,"price":15000}]}}`,
		expectedPath:   "/cart/add",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.AddToCart(ctx, "987654321")
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req AddToCartRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.AssetID != "987654321" {
				t.Errorf("Expected assetId '987654321', got %s", req.AssetID)
			}
		},
		assertResponse: func(t *testing.T, cart CartPayload) {
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
		},
	})

	runAPITest(t, testCase[CartPayload]{
		name:           "AddToCart_error",
		mockStatus:     200,
		mockResponse:   `{"err":true,"success":false,"message":"Item not found","content":null}`,
		expectedPath:   "/cart/add",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.AddToCart(ctx, "invalid_id")
		},
		assertError: func(t *testing.T, err error) {
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
		},
	})
}

func TestBulkAddToCart(t *testing.T) {
	runAPITest(t, testCase[CartPayload]{
		name:           "BulkAddToCart_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"cart":[{"cartId":44,"assetId":"111,222,333","count":3,"item_id":5678,"price":14000}]}}`,
		expectedPath:   "/cart/bulk",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.BulkAddToCart(ctx, 5678, 5, "76561198000000000")
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req BulkAddToCartRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
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
		},
		assertResponse: func(t *testing.T, cart CartPayload) {
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
		},
	})

	runAPITest(t, testCase[CartPayload]{
		name:           "BulkAddToCart_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"message":"Forbidden","content":null}`,
		expectedPath:   "/cart/bulk",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.BulkAddToCart(ctx, 5678, 5, "76561198000000000")
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}

func TestRemoveFromCart(t *testing.T) {
	runAPITest(t, testCase[CartPayload]{
		name:           "RemoveFromCart_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"cart":[]}}`,
		expectedPath:   "/cart/remove",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.RemoveFromCart(ctx, 42)
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req RemoveFromCartRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.CartID != 42 {
				t.Errorf("Expected cartId 42, got %d", req.CartID)
			}
		},
		assertResponse: func(t *testing.T, cart CartPayload) {
			if len(cart.Cart) != 0 {
				t.Fatalf("Expected empty cart, got %d items", len(cart.Cart))
			}
		},
	})

	runAPITest(t, testCase[CartPayload]{
		name:           "RemoveFromCart_internal_error",
		mockStatus:     500,
		mockResponse:   `{"err":true,"success":false,"message":"Internal server error","content":null}`,
		expectedPath:   "/cart/remove",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartPayload, error) {
			return client.RemoveFromCart(ctx, 42)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrInternal) {
				t.Errorf("Expected ErrInternal, got %v", err)
			}
		},
	})
}

func TestUpdateCart(t *testing.T) {
	runAPITest(t, testCase[CartUpdateResponse]{
		name:         "UpdateCart_success",
		mockStatus:   200,
		mockResponse: `{"err":false,"success":true,"message":"","content":{"integrity":{"valid":false,"invalidItems":[{"cartId":42,"assetId":"987654321","reason":"price_changed","currentPrice":16000,"expectedPrice":15000}]},"replaced":[{"cartId":42,"oldAssetId":"987654321","newAssetId":"987654999","reason":"price_changed"}],"removed":[],"cart":[{"cartId":42,"assetId":"987654999","count":1,"item_id":5678,"price":15000}]}}`,
		expectedPath:   "/cart/update",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartUpdateResponse, error) {
			return client.UpdateCart(ctx)
		},
		assertResponse: func(t *testing.T, result CartUpdateResponse) {
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
		},
	})

	runAPITest(t, testCase[CartUpdateResponse]{
		name:           "UpdateCart_unauthorized",
		mockStatus:     401,
		mockResponse:   `{"err":true,"success":false,"message":"Unauthorized","content":null}`,
		expectedPath:   "/cart/update",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (CartUpdateResponse, error) {
			return client.UpdateCart(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}