package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID int64 `json:"user_id"`
	GrandTotal float64 `json:"grand_total"`
}