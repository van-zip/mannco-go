# mannco-go

A Go API client for [Mannco.store](https://mannco.store) — a Team Fortress 2, CS2, Dota 2, Rust, and Steam Community item trading marketplace.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/van-zip/mannco-go.svg)](https://pkg.go.dev/github.com/van-zip/mannco-go)

---

## Features

| Category | Coverage |
|----------|----------|
| **Authentication** | API key → JWT exchange with auto-reauth on 401/403 |
| **Rate Limiting** | Automatic retry with exponential backoff on 429 |
| **Items & Pricing** | Sales graphs, listings, buy orders, pricing (single & bulk up to 100 items) |
| **User Buy Orders** | View your active buy orders (specific item or all) |
| **Buy Orders** | Create, update, remove, and bulk buy orders |
| **User & History** | Balance, transaction / sales / purchase / cashout / balance history, user info, notifications, sales stats/charts, sessions |
| **Inventory** | Items on sale, in inventory, set price, withdraw |
| **Cart & Checkout** | Get, add, bulk add, remove, update cart |
| **Payment** | Create payment sessions (balance & items) |
| **Trading** | *Not yet implemented* |
| **Listing / Deposit** | *Not yet implemented* |

> **Status**: This library covers the endpoints relevant to my own project. Many listing and trading endpoints are not yet implemented. PRs welcome!

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/van-zip/mannco-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpClient := &http.Client{Timeout: 10 * time.Second}
	client := mannco.NewClient("your-mannco-store-api-key", httpClient)

	// Fetch your auth token (JWT)
	fmt.Println("--- Fetching JWT ---")
	err := client.UserLogin(ctx)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Println("--- Fetching User Balance ---")
	balance, err := client.Balance(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v", err)
	}

	fmt.Printf("Current Account Balance: $%.2f\n\n", float64(balance)/100.0)

	fmt.Println("--- Fetching Bulk Pricing Data ---")

	// Max's severed head, Earbuds, Bill's hat ids
	targetIDs := []int{371, 958, 803}
	bulkData, err := client.ItemPricingBulk(ctx, targetIDs)
	if err != nil {
		log.Printf("Warning: Failed bulk pricing call: %v", err)
		return
	}

	fmt.Printf("Successfully analyzed %d items (%d from fresh live updates):\n",
		bulkData.TotalItems,
		bulkData.RefreshedItems,
	)
	for _, item := range bulkData.Items {
		fmt.Printf(" - Item ID %d | Lowest Sale: $%.2f | Suggested Value: $%.2f\n",
			item.ItemID,
			float64(item.Pricing.LowestSalePrice)/100.0,
			float64(item.Pricing.SuggestedPrice)/100.0,
		)
	}
}
```

Run it:

```bash
go run examples/example.go
```

---

## API Reference

### Client

```go
client := mannco.NewClient(apiKey string, httpClient *http.Client)
```

| Method | Description |
|--------|-------------|
| `SetJWT(token string)` | Update the bearer token |
| `GetJWT() string` | Retrieve current token |
| `SetAPIKey(key string)` | Update API key (used for re-auth on 429) |
| `GetAPIKey() string` | Get current API key |
| `SetBaseURL(url string)` | Override API base URL (useful for testing) |
| `GetBaseURL() string` | Get current base URL |

All methods accept `context.Context` as the first argument for cancellation/timeout control.

### Retry Behavior

The client automatically retries requests that return **HTTP 429 (Too Many Requests)** with exponential backoff (1s, 2s, 4s, max 30s) up to 3 attempts. It respects the `Retry-After` header when present.

On **401/403**, the client automatically attempts to re-authenticate using the stored API key before retrying the original request once.

---

## Endpoint Coverage Table

| Endpoint | Method | Tag | Implemented | Function Signature |
|----------|--------|-----|-------------|-------------------|
| `/user/login` | POST | Auth | ✅ | `func (c *Client) UserLogin(ctx context.Context) error` |
| `/item/details/{item}` | GET | Items | ✅ | `func (c *Client) ItemDetails(ctx context.Context, itemID string) (ItemDetailsPayload, error)` |
| `/item/salesGraph/{item}` | GET | Items | ✅ | `func (c *Client) ItemSalesGraph(ctx context.Context, itemID int, period Period) (PriceHistoryPayload, error)` |
| `/item/listing/count/{item}` | GET | Items | ✅ | `func (c *Client) ItemListingCount(ctx context.Context, itemID string, userID string) (ListingCountPayload, error)` |
| `/item/listing/{item}` | GET | Items | ✅ | `func (c *Client) ItemListings(ctx context.Context, itemID int, userID string, opts *ListingOptions) (ListingPayload, error)` |
| `/item/buyorderList/{item}` | GET | Items | ✅ | `func (c *Client) BuyOrderList(ctx context.Context, itemID int) (BuyOrderPayload, error)` |
| `/item/prices` | GET | Items | ✅ | `func (c *Client) ItemPricesByGame(ctx context.Context, game int, outOfStock bool) ([]GameItemPrice, error)` |
| `/item/pricing/{item}` | GET | Items | ✅ | `func (c *Client) ItemPricing(ctx context.Context, itemID int) (PriceItem, error)` |
| `/item/pricing/bulk` | GET | Items | ✅ | `func (c *Client) ItemPricingBulk(ctx context.Context, itemIDs []int) (BulkPricingPayload, error)` |
| `/item/details/fromid/{backpackid}` | GET | Items | ✅ | `func (c *Client) ItemDetailsFromBackpackTF2(ctx context.Context, backpackID string) (BackpackItemPayload, error)` |
| `/item/cs/details/fromid/{backpackid}` | GET | Items | ✅ | `func (c *Client) ItemDetailsFromBackpackCS2(ctx context.Context, backpackID string) (CSBackpackItemPayload, error)` |
| `/offers/received` | GET | Offers | ❌ | — |
| `/offers/my` | GET | Offers | ❌ | — |
| `/offers/create` | POST | Offers | ❌ | — |
| `/offers/accept` | POST | Offers | ❌ | — |
| `/offers/decline` | POST | Offers | ❌ | — |
| `/offers/remove` | POST | Offers | ❌ | — |
| `/item/buyorder` | POST | Buy Orders | ✅ | `func (c *Client) CreateBuyOrder(ctx context.Context, itemID, value, amount int) error` |
| `/item/buyorder/update` | POST | Buy Orders | ✅ | `func (c *Client) UpdateBuyOrder(ctx context.Context, itemID, value, amount int) error` |
| `/item/buyorder/remove` | POST | Buy Orders | ✅ | `func (c *Client) RemoveBuyOrder(ctx context.Context, itemID int) error` |
| `/item/buyorder/bulk` | POST | Buy Orders | ✅ | `func (c *Client) BulkBuyOrders(ctx context.Context, orders []BulkBuyOrderEntry) (BulkBuyOrdersContent, error)` |
| `/user/buyorder/{item}` | GET | Buy Orders | ✅ | `func (c *Client) UserItemBuyOrder(ctx context.Context, itemID int) (UserItemBuyOrderPayload, error)` |
| `/user/getBuyorder` | GET | Buy Orders | ✅ | `func (c *Client) GetUserBuyOrders(ctx context.Context) (UserBuyOrdersPayload, error)` |
| `/payment/{provider}` | POST | Payment | ✅ | `func (c *Client) CreatePayment(ctx context.Context, provider PaymentProvider, req CreatePaymentRequest) (PaymentResponse, error)` |
| `/inventory/onSale` | GET | Inventory | ✅ | `func (c *Client) GetItemsOnSale(ctx context.Context) (InventoryItemsPayload, error)` |
| `/inventory/onInventory` | GET | Inventory | ✅ | `func (c *Client) GetItemsInInventory(ctx context.Context) (InventoryItemsPayload, error)` |
| `/inventory/price` | POST | Inventory | ✅ | `func (c *Client) SetItemPrice(ctx context.Context, ids string, price int) error` |
| `/inventory/withdraw` | POST | Inventory | ✅ | `func (c *Client) WithdrawItems(ctx context.Context, ids string) (WithdrawItemsResponse, error)` |
| `/cart/get` | GET | Cart | ✅ | `func (c *Client) GetCart(ctx context.Context) (CartPayload, error)` |
| `/cart/add` | POST | Cart | ✅ | `func (c *Client) AddToCart(ctx context.Context, assetID string) (CartPayload, error)` |
| `/cart/bulk` | POST | Cart | ✅ | `func (c *Client) BulkAddToCart(ctx context.Context, itemID, count int, sellerUserID string) (CartPayload, error)` |
| `/cart/remove` | POST | Cart | ✅ | `func (c *Client) RemoveFromCart(ctx context.Context, cartID int) (CartPayload, error)` |
| `/cart/update` | POST | Cart | ✅ | `func (c *Client) UpdateCart(ctx context.Context) (CartUpdateResponse, error)` |
| `/deposit/{game}` | GET | Listing | ❌ | — |
| `/deposit/trade` | POST | Listing | ❌ | — |
| `/deposit/instantSell/{game}` | GET | Listing | ❌ | — |
| `/deposit/trade/instant` | POST | Listing | ❌ | — |
| `/deposit/tradeStatus/{tradeid}` | GET | Listing | ❌ | — |
| `/trades/active` | GET | Trading | ❌ | — |
| `/trades/all` | GET | Trading | ❌ | — |
| `/trade/resend` | GET | Trading | ❌ | — |
| `/user/disconnect` | GET | User Account | ✅ | `func (c *Client) Disconnect(ctx context.Context) error` |
| `/user/infos` | GET | User Account | ✅ | `func (c *Client) UserInfo(ctx context.Context) (UserInfoPayload, error)` |
| `/user/balance` | GET | User Account | ✅ | `func (c *Client) Balance(ctx context.Context) (int, error)` |
| `/user/notifications` | GET | User Account | ✅ | `func (c *Client) Notifications(ctx context.Context) (NotificationPayload, error)` |
| `/user/ipList` | GET | User Account | ✅ | `func (c *Client) SessionList(ctx context.Context, page, perPage int, includeExpired bool) (SessionPayload, error)` |
| `/user/store/{identifier}` | GET | User Account | ✅ | `func (c *Client) StoreProfile(ctx context.Context, identifier string) (StoreProfile, error)` |
| `/user/getSalesInfos` | GET | User Account | ✅ | `func (c *Client) SalesStats(ctx context.Context) (json.RawMessage, error)` |
| `/user/getSalesChartInfos` | GET | User Account | ✅ | `func (c *Client) SalesChart(ctx context.Context, period, chart string) (json.RawMessage, error)` |
| `/user/getBalanceHistory` | GET | User Account | ✅ | `func (c *Client) BalanceHistory(ctx context.Context, page, limit int) (BalanceHistoryPayload, error)` |
| `/user/getPurchaseHistory` | GET | User Account | ✅ | `func (c *Client) PurchaseHistory(ctx context.Context, opts *HistoryOptions) (HistoryPayload, error)` |
| `/user/getSalesHistory` | GET | User Account | ✅ | `func (c *Client) SalesHistory(ctx context.Context, opts *HistoryOptions) (HistoryPayload, error)` |
| `/user/getCashoutHistory` | GET | User Account | ✅ | `func (c *Client) CashoutHistory(ctx context.Context, page, limit int) (CashoutHistoryPayload, error)` |
| `/user/getTransactionHistory` | GET | User Account | ✅ | `func (c *Client) TransactionHistory(ctx context.Context, opts *HistoryOptions) (HistoryPayload, error)` |
| `/user/getTransactionDetails` | GET | User Account | ✅ | `func (c *Client) TransactionDetails(ctx context.Context, transactionID string) (TransactionDetailsPayload, error)` |

---

## Key Types

```go
// Pricing periods
type Period string
const (
    Period1Month  Period = "1M"
    Period3Months Period = "3M"
    Period6Months Period = "6M"
    Period1Year   Period = "1Y"
    Period5Years  Period = "5Y"
    PeriodAll     Period = "ALL"
)

// Payment
type PaymentProvider string // "payviox" or "mannco"
type PaymentType string     // "balance" or "items"

type CreatePaymentRequest struct {
    Type           PaymentType
    TOSTimestamp   int64
    Value          int       // for balance
    Items          string    // for items (comma-separated asset IDs)
    PaymentMethod  string
}
```

---

## Testing

```bash
# Unit tests (mock HTTP via httptest)
go test -v ./...

# Integration tests (require MANNCO_API_KEY in .env or env var)
MANNCO_API_KEY=your_key go test -v -tags=integration ./...

# Lint
golangci-lint run
```

---

## License

MIT License — see [LICENSE](LICENSE).

---

PRs welcome for missing endpoints!