package mannco

import (
	"context"
	"encoding/json"
	"testing"
)

func TestBalance(t *testing.T) {
	runAPITest(t, testCase[int]{
		name:           "Balance_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"balance":1250}}`,
		expectedPath:   "/user/balance",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (int, error) {
			return client.Balance(ctx)
		},
		assertResponse: func(t *testing.T, balance int) {
			if balance != 1250 {
				t.Errorf("expected balance 1250, got %d", balance)
			}
		},
	})

	runAPITest(t, testCase[int]{
		name:           "Balance_server_error",
		mockStatus:     500,
		mockResponse:   `{"err":true,"success":false,"message":"Internal server error","content":null}`,
		expectedPath:   "/user/balance",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (int, error) {
			return client.Balance(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if err == nil {
				t.Fatal("expected server error, got nil")
			}
		},
	})

	runAPITest(t, testCase[int]{
		name:           "Balance_malformed_json",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"balance":"not_an_int"}}`,
		expectedPath:   "/user/balance",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (int, error) {
			return client.Balance(ctx)
		},
		assertError: func(t *testing.T, err error) {
			if err == nil {
				t.Fatal("expected JSON decode error, got nil")
			}
		},
	})
}

func TestTransactionHistory(t *testing.T) {
	runAPITest(t, testCase[InventoryPayload]{
		name:           "TransactionHistory_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"count":1,"values":[{"count":1,"date":1704067200,"effect":"","festivized":0,"idbackpack":1,"iditem":440,"image":"","inspect":"","killstreaker":"","level":1,"name":"Test Item","paint":"","parts":"","price":5000,"sheen":"","spell":"","url":""}]}}`,
		expectedPath:   "/user/getTransactionHistory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryPayload, error) {
			return client.TransactionHistory(ctx, &HistoryOptions{Page: 1, Limit: 10})
		},
		assertResponse: func(t *testing.T, res InventoryPayload) {
			if res.Count != 1 {
				t.Errorf("expected count 1, got %d", res.Count)
			}
		},
	})
}

func TestSalesHistory(t *testing.T) {
	runAPITest(t, testCase[InventoryPayload]{
		name:           "SalesHistory_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"count":1,"values":[]}}`,
		expectedPath:   "/user/getSalesHistory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryPayload, error) {
			return client.SalesHistory(ctx, &HistoryOptions{Page: 1, Limit: 5, Period: Period3Months, Search: "test"})
		},
		assertResponse: func(_ *testing.T, _ InventoryPayload) {},
	})
}

func TestPurchaseHistory(t *testing.T) {
	runAPITest(t, testCase[InventoryPayload]{
		name:           "PurchaseHistory_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"count":1,"values":[]}}`,
		expectedPath:   "/user/getPurchaseHistory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (InventoryPayload, error) {
			return client.PurchaseHistory(ctx, &HistoryOptions{Page: 1, Limit: 10})
		},
		assertResponse: func(_ *testing.T, _ InventoryPayload) {},
	})
}

func TestUserInfo(t *testing.T) {
	runAPITest(t, testCase[UserInfoPayload]{
		name:           "UserInfo_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"informations":{"steamId":"76561198000000000","balance":250000,"name":"Player","image":"https://example.com/img.png","tradeurl":"https://steamcommunity.com/tradeoffer/new/?partner=123","2fa":"true","shorturl":"player","notification":1}}}`,
		expectedPath:   "/user/infos",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (UserInfoPayload, error) {
			return client.UserInfo(ctx)
		},
		assertResponse: func(t *testing.T, res UserInfoPayload) {
			if res.Informations.SteamID != "76561198000000000" {
				t.Errorf("expected steamId, got %q", res.Informations.SteamID)
			}
		},
	})
}

func TestNotifications(t *testing.T) {
	runAPITest(t, testCase[NotificationPayload]{
		name:           "Notifications_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"alertCount":0,"messagesCount":0,"offersCount":2,"tradesCount":1}}`,
		expectedPath:   "/user/notifications",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (NotificationPayload, error) {
			return client.Notifications(ctx)
		},
		assertResponse: func(t *testing.T, res NotificationPayload) {
			if res.OffersCount != 2 {
				t.Errorf("expected offersCount 2, got %d", res.OffersCount)
			}
		},
	})
}

func TestBalanceHistory(t *testing.T) {
	runAPITest(t, testCase[BalanceHistoryPayload]{
		name:           "BalanceHistory_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"values":[],"count":0}}`,
		expectedPath:   "/user/getBalanceHistory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (BalanceHistoryPayload, error) {
			return client.BalanceHistory(ctx, 1, 10)
		},
		assertResponse: func(_ *testing.T, _ BalanceHistoryPayload) {},
	})
}

func TestCashoutHistory(t *testing.T) {
	runAPITest(t, testCase[CashoutHistoryPayload]{
		name:           "CashoutHistory_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"values":[]}}`,
		expectedPath:   "/user/getCashoutHistory",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (CashoutHistoryPayload, error) {
			return client.CashoutHistory(ctx, 0, 20)
		},
		assertResponse: func(_ *testing.T, _ CashoutHistoryPayload) {},
	})
}

func TestTransactionDetails(t *testing.T) {
	runAPITest(t, testCase[TransactionDetailsPayload]{
		name:           "TransactionDetails_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"transaction":{}}}`,
		expectedPath:   "/user/getTransactionDetails",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (TransactionDetailsPayload, error) {
			return client.TransactionDetails(ctx, "txn_123")
		},
		assertResponse: func(_ *testing.T, _ TransactionDetailsPayload) {},
	})
}

func TestStoreProfile(t *testing.T) {
	runAPITest(t, testCase[StoreProfile]{
		name:           "StoreProfile_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"steamId":"76561198000000000","name":"Player","image":"https://example.com/img.png"}}`,
		expectedPath:   "/user/store/player",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (StoreProfile, error) {
			return client.StoreProfile(ctx, "player")
		},
		assertResponse: func(t *testing.T, res StoreProfile) {
			if res.Name != "Player" {
				t.Errorf("expected name 'Player', got %q", res.Name)
			}
		},
	})
}

func TestSalesStats(t *testing.T) {
	runAPITest(t, testCase[json.RawMessage]{
		name:           "SalesStats_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{}}`,
		expectedPath:   "/user/getSalesInfos",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (json.RawMessage, error) {
			return client.SalesStats(ctx)
		},
		assertResponse: func(_ *testing.T, _ json.RawMessage) {},
	})
}

func TestSalesChart(t *testing.T) {
	runAPITest(t, testCase[json.RawMessage]{
		name:           "SalesChart_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{}}`,
		expectedPath:   "/user/getSalesChartInfos",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (json.RawMessage, error) {
			return client.SalesChart(ctx, "1M", "1")
		},
		assertResponse: func(_ *testing.T, _ json.RawMessage) {},
	})
}

func TestSessionList(t *testing.T) {
	runAPITest(t, testCase[SessionPayload]{
		name:           "SessionList_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"values":[],"count":0}}`,
		expectedPath:   "/user/ipList",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (SessionPayload, error) {
			return client.SessionList(ctx, 0, 10, true)
		},
		assertResponse: func(_ *testing.T, _ SessionPayload) {},
	})
}

func TestDisconnect(t *testing.T) {
	runAPITest(t, testCase[DisconnectPayload]{
		name:           "Disconnect_success",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"disconnect":true}}`,
		expectedPath:   "/user/disconnect",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (DisconnectPayload, error) {
			err := client.Disconnect(ctx)
			return DisconnectPayload{Disconnect: true}, err
		},
		assertResponse: func(t *testing.T, res DisconnectPayload) {
			if !res.Disconnect {
				t.Error("expected disconnect true")
			}
		},
	})

	runAPITest(t, testCase[DisconnectPayload]{
		name:           "Disconnect_false",
		mockStatus:     200,
		mockResponse:   `{"err":false,"success":true,"message":"","content":{"disconnect":false}}`,
		expectedPath:   "/user/disconnect",
		expectedMethod: "GET",
		runTest: func(ctx context.Context, client *Client) (DisconnectPayload, error) {
			err := client.Disconnect(ctx)
			return DisconnectPayload{}, err
		},
		assertError: func(t *testing.T, err error) {
			if err == nil {
				t.Fatal("expected error for disconnect false, got nil")
			}
		},
	})
}
