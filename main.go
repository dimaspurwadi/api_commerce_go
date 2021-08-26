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
		productWarehouse.PUT("/restock", middleware.CheckJWT(1), core.RestockProductWarehouse)
		productWarehouse.DELETE("/delete/:sku", middleware.CheckJWT(1), core.DeleteProductWarehouse)
	}

	productDisplay := r.Group("/display")
	{
		productDisplay.GET("/", core.GetProductDisplay)
		productDisplay.GET("/:sku", core.GetProductDisplayBySku)
		productDisplay.POST("/insert", middleware.CheckJWT(1), core.InsertProductDisplay)
		productDisplay.PUT("/update", middleware.CheckJWT(1), core.UpdateProductDisplay)
		productDisplay.PUT("/restock", middleware.CheckJWT(1), core.RestockProductDisplay)
		productDisplay.DELETE("/delete/:sku", middleware.CheckJWT(1), core.DeleteProductDisplay)
	}

	cart := r.Group("/cart")
	{
		cart.GET("/", middleware.CheckJWT(0), core.GetCart) 
		cart.POST("/add", middleware.CheckJWT(0), core.AddToCart)
		cart.DELETE("/delete/:sku", middleware.CheckJWT(0), core.DeleteProductFromCart)
		cart.POST("/checkout", middleware.CheckJWT(0), core.CheckoutCart)
	}

	transaction := r.Group("/transaction")
	{
		transaction.GET("/", middleware.CheckJWT(0), core.GetTransactionHistory) 
	}

	report := r.Group("/report")
	{
		report.GET("/", middleware.CheckJWT(1), core.GetAllReport) 
	}

	r.Run(":" + os.Getenv("API_PORT"))
}