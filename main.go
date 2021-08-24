package main

import "github.com/subosito/gotenv"
import "github.com/gin-gonic/gin"
import "os"
import "Course-Net/final_project2/config"
import "Course-Net/final_project2/core"
import "Course-Net/final_project2/middleware"

func main() {
	gotenv.Load()
	config.InitDB()

	sqlDB, err := config.Db.DB()
	if err != nil {
		panic("Failed to get generic database object")
	}
	defer sqlDB.Close()

	r := gin.Default()
	auth := r.Group("/auth")
	{
		auth.POST("/register", core.Register)
		auth.POST("/login", core.Login)
		auth.POST("/logout", middleware.CheckJWTWithoutRole(), core.Logout)
	}

	productWarehouse := r.Group("/stock")
	{
		productWarehouse.GET("/", middleware.CheckJWT(1), core.GetProductWarehouse)
		productWarehouse.GET("/:sku", middleware.CheckJWT(1), core.GetProductWarehouseBySku)
		productWarehouse.POST("/insert", middleware.CheckJWT(1), core.InsertProductWarehouse)
		productWarehouse.PUT("/update", middleware.CheckJWT(1), core.UpdateProductWarehouse)
	}

	r.Run(":" + os.Getenv("API_PORT"))
}