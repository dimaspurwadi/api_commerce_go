package config

import model "Course-Net/final_project2/models"
import "os"
import "gorm.io/driver/mysql"
import "gorm.io/gorm"

var Db *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	Db.AutoMigrate(&model.User{})
}

