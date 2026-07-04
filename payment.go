package mannco

import (
	"context"
	"encoding/json"
	"fmt"
)

// PaymentProvider is the payment gateway to use
type PaymentProvider string

// Valid payment providers
const (
	ProviderPayviox PaymentProvider = "payviox"
	ProviderMannco  PaymentProvider = "mannco"
)

// PaymentType is the kind of payment being created
type PaymentType string

// Valid payment types
const (
	PaymentTypeBalance PaymentType = "balance"
	PaymentTypeItems   PaymentType = "items"
)

// CreatePaymentRequest for POST /payment/{provider}
type CreatePaymentRequest struct {
	Type          PaymentType `json:"type"`
	TOSTimestamp  int64       `json:"tos_timestamp"`
	Value         int         `json:"value,omitempty"`
	Items         string      `json:"items,omitempty"`
	PaymentMethod string      `json:"payment_method,omitempty"`
}

// PaymentResponse is the response from a payment initiation
type PaymentResponse struct {
	URL     string `json:"url,omitempty"`
	Message string `json:"message,omitempty"`
}

// CreatePayment initiates a payment session with the given provider
func (c *Client) CreatePayment(ctx context.Context, provider PaymentProvider, req CreatePaymentRequest) (PaymentResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return PaymentResponse{}, fmt.Errorf("%w: error encoding json for create payment: %w", ErrInternal, err)
	}
	return executeRequest[PaymentResponse](ctx, c, "POST", "payment/"+string(provider), jsonData, nil)
}

