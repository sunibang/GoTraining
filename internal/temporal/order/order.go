package order

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID        uuid.UUID  `json:"id"`
	LineItems []LineItem `json:"line_items"`
}

type LineItem struct {
	ProductID    uuid.UUID       `json:"product_id"`
	Quantity     int32           `json:"quantity"`
	PricePerItem decimal.Decimal `json:"price_per_item"`
}
