package model

import "time"

type StockMovement struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	LocationID int       `json:"location_id"`
	Type       string    `json:"type"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
}
