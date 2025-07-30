package models

import "time"

type Payment struct {
	Billing struct {
		City    string `json:"city"`
		Country string `json:"country"`
		State   string `json:"state"`
		Street  string `json:"street"`
		Zipcode string `json:"zipcode"`
	} `json:"billing"`

	BrandID            string    `json:"brand_id"`
	BusinessID         string    `json:"business_id"`
	CardIssuingCountry *string   `json:"card_issuing_country"`
	CardLastFour       string    `json:"card_last_four"`
	CardNetwork        string    `json:"card_network"`
	CardType           string    `json:"card_type"`
	CreatedAt          time.Time `json:"created_at"`
	Currency           string    `json:"currency"`

	Customer struct {
		CustomerID string `json:"customer_id"`
		Email      string `json:"email"`
		Name       string `json:"name"`
	} `json:"customer"`

	DigitalProductsDelivered bool   `json:"digital_products_delivered"`
	DiscountID               string `json:"discount_id"`

	Disputes []struct {
		Amount        string    `json:"amount"`
		BusinessID    string    `json:"business_id"`
		CreatedAt     time.Time `json:"created_at"`
		Currency      string    `json:"currency"`
		DisputeID     string    `json:"dispute_id"`
		DisputeStage  string    `json:"dispute_stage"`
		DisputeStatus string    `json:"dispute_status"`
		PaymentID     string    `json:"payment_id"`
		Remarks       string    `json:"remarks"`
	} `json:"disputes"`

	ErrorCode         string                 `json:"error_code"`
	ErrorMessage      string                 `json:"error_message"`
	Metadata          map[string]interface{} `json:"metadata"`
	PaymentID         string                 `json:"payment_id"`
	PaymentLink       string                 `json:"payment_link"`
	PaymentMethod     string                 `json:"payment_method"`
	PaymentMethodType string                 `json:"payment_method_type"`

	ProductCart []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"product_cart"`

	Refunds []struct {
		Amount     int       `json:"amount"`
		BusinessID string    `json:"business_id"`
		CreatedAt  time.Time `json:"created_at"`
		Currency   *string   `json:"currency"`
		IsPartial  bool      `json:"is_partial"`
		PaymentID  string    `json:"payment_id"`
		Reason     string    `json:"reason"`
		RefundID   string    `json:"refund_id"`
		Status     string    `json:"status"`
	} `json:"refunds"`

	SettlementAmount   int       `json:"settlement_amount"`
	SettlementCurrency string    `json:"settlement_currency"`
	SettlementTax      int       `json:"settlement_tax"`
	Status             *string   `json:"status"`
	SubscriptionID     string    `json:"subscription_id"`
	Tax                int       `json:"tax"`
	TotalAmount        int       `json:"total_amount"`
	UpdatedAt          time.Time `json:"updated_at"`
}
