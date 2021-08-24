package models

import "gorm.io/gorm"

type ProductWarehouse struct {
	gorm.Model
	Sku string `gorm:"type:varchar(50);unique" json:"sku"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Price float64  `json:"price"`
	Qty int64	`json:"qty"`
}

type ProductWarehouseRequest struct {
	Sku string `json:"sku" binding:"required"`
	Name string `json:"name" binding:"required"`
	Price float64  `json:"price" binding:"required"`
	Qty int64	`json:"qty" binding:"required"`
	
}
