package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uint64      `json:"order_id"`
	CustomerID  uuid.UUID   `json:"customer_id"`
	LineItems   []LineItems `json:"line_items"`
	CreatedAt   time.Time   `json:"created_at"`
	ShippedAT   time.Time   `json:"shipped_at"`
	CompletedAt time.Time   `json:"completed_at"`
}

type LineItems struct {
	ItemID   uuid.UUID `json:"item_id"`
	Price    uint      `json:"price"`
	Quantity uint      `json:"quantity"`
}
