package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
import "strconv"
import "fmt"
import "time"

func GetAllReport(c *gin.Context) {
	time := time.Now()
	year := strconv.Itoa(time.Year())
	month := int(time.Month())
	var monthString string
	if (month < 10) {
		monthString = fmt.Sprintf("0%s",strconv.Itoa(month))
	}
	day := strconv.Itoa(time.Day())
	fmt.Println(year)
	fmt.Println(monthString)
	fmt.Println(day)
	filter,  _ := strconv.ParseBool(c.Query("filter"))
	dayStart := c.Query("dayStart")
	dayEnd := c.Query("dayEnd")
	monthStart := c.Query("monthStart")
	monthEnd := c.Query("monthEnd")
	yearStart := c.Query("yearStart")
	yearEnd := c.Query("yearEnd")
	
	var reportTransaction []model.Transaction
	err := config.Db.Find(&reportTransaction, "1 = ?", 1)
	if filter {
		if dayStart == "" && 
		   dayEnd == "" &&
		   monthStart == "" &&
		   monthEnd == "" &&
		   yearStart == "" &&
		   yearEnd == "" {
		   c.JSON(http.StatusBadRequest, gin.H{
			   "status" : "Bad Request",
			   "messages" : "params dayStart, dayEnd, monthStart, monthEnd, yearStart, and yearEnd is required",
		   })
		   c.Abort()
		   return
		} 
		dateStart := yearStart+"-"+monthStart+"-"+dayStart
		dateEnd := yearEnd+"-"+monthEnd+"-"+dayEnd
		err = config.Db.Find(&reportTransaction, "1 = ? and date(created_at) >= ? and date(created_at) <= ? ", 1, dateStart, dateEnd)
	} 	
	
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	var transactionDetail *model.TransactionDetail
 	var transactionsObj []map[string]interface{}
	for _, transaction := range reportTransaction {
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