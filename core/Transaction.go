package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
// import "reflect"
import "fmt"
// import "io/ioutil"

func GetTransactionHistory(c *gin.Context) {
	var transactionDetail []model.TransactionDetail
	var transactions []model.Transaction
	userID := ClaimToken(c)
	err := config.Db.Find(&transactions, "user_id = ?", userID)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
	} else {
		var transactionsObj []map[string]interface{}
		for _, transaction := range transactions {
			config.Db.Find(&transactionDetail, "transaction_id = ?", transaction.ID)
			var transactionList = map[string]interface{}{
				"ID":          	 transaction.ID,
				"CreatedAt":     transaction.CreatedAt,
				"UpdatedAt":     transaction.UpdatedAt,
				"DeletedAt":     transaction.DeletedAt,
				"User_id":       transaction.UserID,
				"SubTotal": 	 transaction.SubTotal,
				"DiscountTotal": 	 transaction.DiscountTotal,
				"GrandTotal": 	 transaction.GrandTotal,
				"detail" : transactionDetail, 
			}
			transactionsObj = append(transactionsObj, transactionList)
		}
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : transactionsObj,
		})
	}
}