package core

import "Course-Net/final_project2/config"
import model "Course-Net/final_project2/models"
import "net/http"
import "github.com/gin-gonic/gin"
import "strconv"
import "fmt"
import "time"

var ReportTransactionDetail []model.TransactionDetail
var ReportTransaction []model.Transaction

func GetAllReport(c *gin.Context) {
	userID := ClaimToken(c)
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
	
	// sql := fmt.Sprintf("SELECT * FROM transactions WHERE user_id = %d", userID)
	// sqlWhere := sql
	err := config.Db.Find(&ReportTransaction, "user_id = ?", userID)
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
		// sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		err = config.Db.Find(&ReportTransaction, "user_id = ? and date(created_at) >= ? and date(created_at) <= ? ", userID, dateStart, dateEnd)
	} 	
	// fmt.Println(sqlWhere)
	// err := config.Db.Raw(sqlWhere).Scan(&ReportTransaction)
	
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"status" : "Internal Server Error",
			"messages" : fmt.Sprintf("%s", err.Error),
		})
		c.Abort()
		return
	}
	var transactions[]map[string]interface{}
	for _, transaction := range ReportTransaction {
		config.Db.Find(&TransactionDetail, "transaction_id = ?", transaction.ID)
		var transactionList = map[string]interface{}{
			"ID":          	 transaction.ID,
			"CreatedAt":     transaction.CreatedAt,
			"UpdatedAt":     transaction.UpdatedAt,
			"DeletedAt":     transaction.DeletedAt,
			"User_id":       transaction.UserID,
			"SubTotal": 	 transaction.SubTotal,
			"DiscountTotal": 	 transaction.DiscountTotal,
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

// if dayStart != "" && 
		//    dayEnd == "" &&
		//    monthStart == "" &&
		//    monthEnd == "" &&
		//    yearStart == "" &&
		//    yearEnd == "" {
		// 	dateStart := year+"-"+monthString+"-"+dayStart
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s'", sql, dateStart)
		// } else if dayStart != "" && 
		// 		  dayEnd != "" &&
		// 	  	  monthStart == "" &&
		// 		  monthEnd == "" &&
		// 		  yearStart == "" &&
		// 		  yearEnd == "" {
		// 	dateStart := year+"-"+monthString+"-"+dayStart
		// 	dateEnd := year+"-"+monthString+"-"+dayEnd
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		// } else if dayStart != "" && 
		// 		  dayEnd != "" &&
		// 		  monthStart != "" &&
		// 		  monthEnd == "" &&
		// 		  yearStart == "" &&
		// 		  yearEnd == "" {
		// 	dateStart := year+"-"+monthStart+"-"+dayStart
		// 	dateEnd := year+"-"+monthString+"-"+dayEnd
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		// } else if dayStart != "" && 
		// 		  dayEnd != "" &&
		// 		  monthStart != "" &&
		// 		  monthEnd != "" &&
		// 		  yearStart == "" &&
		// 		  yearEnd == "" {
		// 	dateStart := year+"-"+monthStart+"-"+dayStart
		// 	dateEnd := year+"-"+monthEnd+"-"+dayEnd
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		// } else if dayStart != "" && 
		// 		  dayEnd != "" &&
		// 		  monthStart != "" &&
		// 		  monthEnd != "" &&
		// 		  yearStart != "" &&
		// 		  yearEnd == "" {
		// 	dateStart := yearStart+"-"+monthStart+"-"+dayStart
		// 	dateEnd := year+"-"+monthEnd+"-"+dayEnd
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		// } else if dayStart != "" && 
		// 		  dayEnd != "" &&
		// 		  monthStart != "" &&
		// 		  monthEnd != "" &&
		// 		  yearStart != "" &&
		// 		  yearEnd != "" {
		// 	dateStart := yearStart+"-"+monthStart+"-"+dayStart
		// 	dateEnd := yearEnd+"-"+monthEnd+"-"+dayEnd
		// 	sqlWhere = fmt.Sprintf("%s and date(created_at) >= '%s' and date(created_at) <= '%s' ", sql, dateStart, dateEnd)
		// }