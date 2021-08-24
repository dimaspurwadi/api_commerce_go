package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
// import "reflect"
import "fmt"
// import "io/ioutil"

var TransactionDetail []model.TransactionDetail
var Transactions []model.Transaction

func GetTransactionHistory(c *gin.Context) {
	userID := ClaimToken(c)
	err := config.Db.Find(&Transactions, "user_id = ?", userID)

	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
	} else {
		var transactions []map[string]interface{}
		for _, transaction := range Transactions {
			config.Db.Find(&TransactionDetail, "transaction_id = ?", transaction.ID)
			var transactionList = map[string]interface{}{
				"ID":          	 transaction.ID,
				"CreatedAt":     transaction.CreatedAt,
				"UpdatedAt":     transaction.UpdatedAt,
				"DeletedAt":     transaction.DeletedAt,
				"User_id":       transaction.UserID,
				"GrandTotal": 	 transaction.GrandTotal,
				"detail" : TransactionDetail, 
			}
			transactions = append(transactions, transactionList)
		}
		c.JSON(http.StatusOK, gin.H {
			"status" : "success",
			"messages" : "Success Get Product",
			"data" : transactions,
		})
	}
}