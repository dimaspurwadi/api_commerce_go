package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
import "fmt"

var Carts []model.Cart
var Cart model.Cart

func GetCart(c *gin.Context) {
	userID := ClaimToken(c)
	err := config.Db.Find(&Carts, "user_id = ?", userID)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : "Product Not Found",
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Cart",
			"data" : Carts,
		})
	}
}

func AddToCart(c *gin.Context) {
	userID := ClaimToken(c)
	
	var user model.User
	err := config.Db.First(&user, "id = ?", userID)
	
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status" : "Not Found",
			"messages" : "User Not Found",
		})
		c.Abort()
		return
	}
	
	var json model.CartRequest
	
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status" : "Bad Request",
			"messages" : fmt.Sprintf("%s", err.Error()),
		})
		c.Abort()
		return
	}
	
	var productDisplay model.ProductDisplay
	err = config.Db.First(&productDisplay, "sku = ?", json.Sku)
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : "Data Sku "+json.Sku+" at ProductDisplay  Not Found",
		})
		c.Abort()
		return
	}

	var cart model.Cart
	err = config.Db.First(&cart, "user_id = ? and sku = ?", userID, json.Sku)

	if err.Error != nil {
		cart = model.Cart{
			UserID: int64(userID),
			Sku: json.Sku,
			Qty: json.Qty,
		}
		err = config.Db.Create(&cart)
		
	} else {
		cart.Qty += json.Qty
		err = config.Db.Save(&cart)
	}

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
		"messages" : "Product Cart success insert",
		"data" : cart,
	})
}

func DeleteProductFromCart(c *gin.Context) {
	userID := ClaimToken(c)
	err := config.Db.First(&Cart,"sku = ? and user_id", c.Param("sku"), userID)

	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status" : "Not Found",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	} 

	err = config.Db.Delete(&Cart)
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
		"messages" : "Success Delete Product From Cart",
		"data" : Cart,
	})
}

func CheckoutCart(c *gin.Context) {
	userID := ClaimToken(c)
	err := config.Db.Find(&Carts, "user_id = ?", userID)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : "Product Not Found",
		})
	}
	/*check Qty available on display*/
	var subTotal float64
	for _, cart := range Carts {
		statusQty := validateCart(cart)
		if !statusQty {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status" : "Unprocessable",
				"messages" : fmt.Sprintf("Qty of sku %s not enough", cart.Sku),
			})
			c.Abort()
			return
		}
		productDisplay := getProductDisplayBySku(cart.Sku)
		subTotal += productDisplay.Price * float64(cart.Qty)
	}

	discountTotal := 0.00

	transaction := model.Transaction{
		UserID : int64(userID),
		SubTotal: subTotal,
		DiscountTotal : discountTotal,
		GrandTotal : subTotal + discountTotal,
	}
	
	err = config.Db.Create(&transaction)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status" : "Internal Server Error",
			"messages" : "Failed to Save Transaction",
		})
		c.Abort()
		return	
	}

	for _, cart := range Carts {
		productDisplay := getProductDisplayBySku(cart.Sku)
		Total := productDisplay.Price * float64(cart.Qty)
		transactionDetail := model.TransactionDetail{
			TransactionID: int64(transaction.ID),
			Sku: cart.Sku,
			Qty: cart.Qty,
			Price: productDisplay.Price,
			Total: Total,
		}
		err := config.Db.Create(&transactionDetail)
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status" : "Internal Server Error",
				"messages" : fmt.Sprintf("%s", err.Error),
			})
			c.Abort()
			return		
		}
		upDateQuantityOnDisplay(cart)
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : "success",
		"messages" : "checkout Cart Sukses",
		"data" : transaction,
	})
	c.Abort()
	return
}

func getProductWarehouseBySku(sku string) *model.ProductWarehouse {
	var productByWarehouse model.ProductWarehouse
	err := config.Db.First(&productByWarehouse, "sku = ?", sku)
	if err.Error != nil {
		return nil
	}
	return &productByWarehouse
}

func getProductDisplayBySku(sku string) *model.ProductDisplay {
	var productDisplay model.ProductDisplay
	err := config.Db.First(&productDisplay, "sku = ?", sku)
	if err.Error != nil {
		return nil
	}
	return &productDisplay
}

func validateCart(cart model.Cart) bool {
	var productDisplay model.ProductDisplay
	err := config.Db.First(&productDisplay, "sku = ?", cart.Sku)
	if err.Error != nil {
		return false
	}
	if (cart.Qty > productDisplay.Qty) {
		return false
	}
	return true
}

func upDateQuantityOnDisplay(cart model.Cart) {
	var productDisplay model.ProductDisplay
	config.Db.First(&productDisplay, "sku = ?", cart.Sku)
	productDisplay.Qty -= cart.Qty
	config.Db.Save(&productDisplay)
}