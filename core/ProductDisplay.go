package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
// import "reflect"
import "fmt"
// import "io/ioutil"

var ProductDisplays []model.ProductDisplay
var ProductDisplay model.ProductDisplay

func GetProductDisplay(c *gin.Context) {
	err := config.Db.Find(&ProductDisplays)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : "Product Not Found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : ProductDisplays,
		})
	}
}

func InsertProductDisplay(c *gin.Context) {
	var json model.ProductDisplayRequest
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
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : "Data at ProductWarehouse Not Found",
		})
		c.Abort()
		return
	}

	if (productWarehouse.Qty < json.Qty) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status" : "Unprocessable",
			"messages" : "Total Qty di Product Display melebihi total Qty di Product Warehouse pada SKU " + json.Sku,
		})
		c.Abort()
		return
	}

	productDisplay := model.ProductDisplay{
		Sku : json.Sku,
		Name : productWarehouse.Name,
		Price : json.Price,
		Qty : json.Qty,
	}
	err = config.Db.Create(&productDisplay)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"message" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	} 
	
	productWarehouse.Qty -= json.Qty
	err = config.Db.Save(&productWarehouse)

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
		"messages" : "Product Display success insert",
		"data" : productDisplay,
	})
}

func GetProductDisplayBySku(c *gin.Context) {
	err := config.Db.First(&ProductDisplay,"sku = ?", c.Param("sku"))

	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : "Product Not Found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : ProductDisplay,
		})
	}
}

func UpdateProductDisplay(c *gin.Context) {
	var json model.UpdateProductDisplayRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Bad Request",
			"messages" : fmt.Sprintf("%s", err.Error()),
		})
		c.Abort()
		return
	}
	var productDisplay model.ProductDisplay
	err := config.Db.First(&productDisplay, "sku = ?", json.Sku)
	if err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Not Found",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	productDisplay.Name = json.Name
	productDisplay.Price = json.Price
	productDisplay.Qty = json.Qty

	result := config.Db.Save(&productDisplay)
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
		"data" : productDisplay,
	})
}

func RestockProductDisplay(c *gin.Context) {
	var json model.RestockProductDisplayRequest
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
			"status" : "Not Found",
			"messages" : "Product SKU "+json.Sku+" Not Found",
		})
		c.Abort()
		return
	}
	
	productDisplay := model.ProductDisplay{
		Sku : json.Sku,
		Name : productWarehouse.Name,
		Price : productWarehouse.Price,
		Qty : productWarehouse.Qty,
	}

	err = config.Db.Create(&productDisplay)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"message" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	} 
	
	productWarehouse.Qty -= json.Qty
	err = config.Db.Save(&productWarehouse)

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
		"messages" : "Product Display Restock Success",
		"data" : productDisplay,
	})
}

func DeleteProductDisplay(c *gin.Context) {
	err := config.Db.First(&ProductWarehouse,"sku = ?", c.Param("sku"))

	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	} 

	var productDisplay model.ProductDisplay
	err = config.Db.First(&productDisplay, "sku = ?", c.Param("sku"))

	ProductWarehouse.Qty += productDisplay.Qty
	err = config.Db.Save(&ProductWarehouse)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	
	err = config.Db.Delete(&ProductWarehouse)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}

	err = config.Db.Delete(&productDisplay)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	
	c.JSON(http.StatusOK, gin.H {
		"status" : "success",
		"messages" : "Success Delete Product",
		"data" : ProductDisplay,
	})
}