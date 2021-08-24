package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID int64 `json:"user_id"`
	Sku string `json:"sku"`
	Qty int64	`json:"qty"`
}

type CartRequest struct {
	Sku string `json:"sku" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`	
}

type CheckoutRequest struct {
	Sku string `json:"sku" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`	
}