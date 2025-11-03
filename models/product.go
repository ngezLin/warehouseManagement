package model

type Product struct {
	ID       int    `json:"id"`
	SKUName  string `json:"sku_name"`
	Quantity int    `json:"quantity"`
}