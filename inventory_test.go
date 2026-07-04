package mannco

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

func TestGetItemsOnSale(t *testing.T) {
	runAPITest(t, testCase[InventoryItemsPayload]{
		name:         "GetItemsOnSale_success",
		mockStatus:   200,
		mockResponse: `{"err":false,"success":true,"message":"","content":{"items":[{"ids":"987654321,987654322","count":2,"item_id":5678,"state":1,"price":15000,"name":"Unusual Burning Flames Team Captain","game":440}],"count":25}}`,
		expectedPath:   "/inventory/onSale",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryItemsPayload, error) {
			return client.GetItemsOnSale(ctx)
		},
		assertResponse: func(t *testing.T, payload InventoryItemsPayload) {
			if payload.Count != 25 {
				t.Errorf("Expected count 25, got %d", payload.Count)
			}
			if len(payload.Items) != 1 {
				t.Fatalf("Expected 1 item, got %d", len(payload.Items))
			}
			item := payload.Items[0]
			if item.IDs != "987654321,987654322" {
				t.Errorf("Expected ids '987654321,987654322', got %s", item.IDs)
			}
			if item.Count != 2 {
				t.Errorf("Expected count 2, got %d", item.Count)
			}
			if item.ItemID != 5678 {
				t.Errorf("Expected item_id 5678, got %d", item.ItemID)
			}
			if item.State != 1 {
				t.Errorf("Expected state 1, got %d", item.State)
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

	runAPITest(t, testCase[InventoryItemsPayload]{
		name:           "GetItemsOnSale_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"content":"forbidden"}`,
		expectedPath:   "/inventory/onSale",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryItemsPayload, error) {
			return client.GetItemsOnSale(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}

func TestGetItemsInInventory(t *testing.T) {
	runAPITest(t, testCase[InventoryItemsPayload]{
		name:         "GetItemsInInventory_success",
		mockStatus:   200,
		mockResponse: `{"err":false,"success":true,"message":"","content":{"items":[{"ids":"987654322","count":1,"item_id":5679,"state":0,"name":"Strange Shotgun","game":440}],"count":10}}`,
		expectedPath:   "/inventory/onInventory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryItemsPayload, error) {
			return client.GetItemsInInventory(ctx)
		},
		assertResponse: func(t *testing.T, payload InventoryItemsPayload) {
			if payload.Count != 10 {
				t.Errorf("Expected count 10, got %d", payload.Count)
			}
			if len(payload.Items) != 1 {
				t.Fatalf("Expected 1 item, got %d", len(payload.Items))
			}
			item := payload.Items[0]
			if item.IDs != "987654322" {
				t.Errorf("Expected ids '987654322', got %s", item.IDs)
			}
			if item.Count != 1 {
				t.Errorf("Expected count 1, got %d", item.Count)
			}
			if item.ItemID != 5679 {
				t.Errorf("Expected item_id 5679, got %d", item.ItemID)
			}
			if item.State != 0 {
				t.Errorf("Expected state 0, got %d", item.State)
			}
			if item.Price != 0 {
				t.Errorf("Expected price 0 (not on sale), got %d", item.Price)
			}
			if item.Name != "Strange Shotgun" {
				t.Errorf("Expected name 'Strange Shotgun', got %s", item.Name)
			}
			if item.Game != 440 {
				t.Errorf("Expected game 440, got %d", item.Game)
			}
		},
	})

	runAPITest(t, testCase[InventoryItemsPayload]{
		name:           "GetItemsInInventory_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"content":"forbidden"}`,
		expectedPath:   "/inventory/onInventory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryItemsPayload, error) {
			return client.GetItemsInInventory(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}

func TestSetItemPrice(t *testing.T) {
	runAPITest(t, testCase[json.RawMessage]{
		name:           "SetItemPrice_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"message":"ok"}}`,
		expectedPath:   "/inventory/price",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (json.RawMessage, error) {
			return nil, client.SetItemPrice(ctx, "12345678,87654321", 150)
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req struct {
				IDs   string `json:"ids"`
				Price int    `json:"price"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.IDs != "12345678,87654321" {
				t.Errorf("Expected ids '12345678,87654321', got %s", req.IDs)
			}
			if req.Price != 150 {
				t.Errorf("Expected price 150, got %d", req.Price)
			}
		},
		assertResponse: func(_ *testing.T, _ json.RawMessage) {},
	})

	runAPITest(t, testCase[json.RawMessage]{
		name:           "SetItemPrice_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"content":"forbidden"}`,
		expectedPath:   "/inventory/price",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (json.RawMessage, error) {
			return nil, client.SetItemPrice(ctx, "12345678", 150)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})

	runAPITest(t, testCase[json.RawMessage]{
		name:           "SetItemPrice_business_error",
		mockStatus:     300,
		mockResponse:   `{"err":true,"success":false,"content":"Item not found in inventory"}`,
		expectedPath:   "/inventory/price",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (json.RawMessage, error) {
			return nil, client.SetItemPrice(ctx, "99999999", 150)
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrInternal) {
				t.Errorf("Expected ErrInternal, got %v", err)
			}
		},
	})
}

func TestWithdrawItems(t *testing.T) {
	runAPITest(t, testCase[WithdrawItemsResponse]{
		name:         "WithdrawItems_success",
		mockStatus:   200,
		mockResponse: `{"err":false,"success":true,"message":"","content":{"message":"Items withdrawal processed","updated":3,"locked":0}}`,
		expectedPath:   "/inventory/withdraw",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (WithdrawItemsResponse, error) {
			return client.WithdrawItems(ctx, "12345678,87654321")
		},
		assertRequestBody: func(t *testing.T, body []byte) {
			var req struct {
				IDs string `json:"ids"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
			}
			if req.IDs != "12345678,87654321" {
				t.Errorf("Expected ids '12345678,87654321', got %s", req.IDs)
			}
		},
		assertResponse: func(t *testing.T, resp WithdrawItemsResponse) {
			if resp.Message != "Items withdrawal processed" {
				t.Errorf("Expected message 'Items withdrawal processed', got %s", resp.Message)
			}
			if resp.Updated != 3 {
				t.Errorf("Expected updated 3, got %d", resp.Updated)
			}
			if resp.Locked != 0 {
				t.Errorf("Expected locked 0, got %d", resp.Locked)
			}
		},
	})

	runAPITest(t, testCase[WithdrawItemsResponse]{
		name:           "WithdrawItems_unauthorized",
		mockStatus:     403,
		mockResponse:   `{"err":true,"success":false,"content":"forbidden"}`,
		expectedPath:   "/inventory/withdraw",
		expectedMethod: "POST",
		runTest: func(ctx context.Context, client *Client) (WithdrawItemsResponse, error) {
			return client.WithdrawItems(ctx, "12345678")
		},
		assertError: func(t *testing.T, err error) {
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("Expected ErrUnauthorized, got %v", err)
			}
		},
	})
}