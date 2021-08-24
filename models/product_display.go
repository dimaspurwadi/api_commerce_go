package models

import "gorm.io/gorm"

type ProductDisplay struct {
	gorm.Model
	Sku string `gorm:"type:varchar(50);unique" json:"sku"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Price float64  `json:"price"`
	Qty int64	`json:"qty"`
}

type ProductDisplayRequest struct {
	Sku string `json:"sku" binding:"required"`
	Price float64  `json:"price" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`	
}

type RestockProductDisplayRequest struct {
	Sku string `json:"sku" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`	
}

type UpdateProductDisplayRequest struct {
	Sku string `json:"sku" binding:"required"`
	Name string `json:"name" binding:"required"`
	Price float64  `json:"price" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`	
}