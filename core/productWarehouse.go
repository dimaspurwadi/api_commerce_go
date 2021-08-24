package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
// import "reflect"
import "fmt"
// import "io/ioutil"

var ProductWarehouses []model.ProductWarehouse
var ProductWarehouse model.ProductWarehouse

func GetProductWarehouse(c *gin.Context) {
	err := config.Db.Find(&ProductWarehouses)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : "Product Not Found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : ProductWarehouse,
		})
	}
}

func InsertProductWarehouse(c *gin.Context) {
	var json model.ProductWarehouseRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Bad Request",
			"messages" : fmt.Sprintf("%s", err.Error()),
		})
		c.Abort()
		return
	}
	productWareHouse := model.ProductWarehouse{
		Sku : json.Sku,
		Name : json.Name,
		Price : json.Price,
		Qty : json.Qty,
	}
	err := config.Db.Create(&productWareHouse)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"message" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status" : "success",
		"messages" : "Product Warehouse success insert",
		"data" : productWareHouse,
	})
}

func GetProductWarehouseBySku(c *gin.Context) {
	err := config.Db.First(&ProductWarehouse,"sku = ?", c.Param("sku"))

	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : "Product Not Found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : ProductWarehouse,
		})
	}
}

func UpdateProductWarehouse(c *gin.Context) {
	var json model.ProductWarehouseRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Bad Request",
			"messages" : fmt.Sprintf("%s", err.Error()),
		})
		c.Abort()
		return
	}
	var productWarehouse model.ProductWarehouse
	err := config.Db.First(&productWarehouse, "sku = ?", json.Sku)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Bad Request",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	productWarehouse.Sku = json.Sku
	productWarehouse.Name = json.Name
	productWarehouse.Price = json.Price
	productWarehouse.Qty = json.Qty

	result := config.Db.Save(&productWarehouse)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", result.Error),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status" : "success",
		"messages" : "Success Update Data",
		"data" : productWarehouse,
	})
}