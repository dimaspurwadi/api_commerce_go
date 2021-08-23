package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(150);unique"`
	Fullname string `gorm:"type:varchar(255)"`
	Email string	`gorm:"type:varchar(255)"`
	Address string	
	Password string	`gorm:"type:varchar(255)"`
	Role int64
	Token string 
}
