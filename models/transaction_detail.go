package models

import "gorm.io/gorm"

type TransactionDetail struct {
	gorm.Model
	TransactionID int64 `json:"transaction_id"`
	Sku string `json:"sku"`
	Qty int64 `json:"qty"`
	SubTotal float64 `json:"sub_total"`
}