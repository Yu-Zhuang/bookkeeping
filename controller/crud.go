package controller

import (
	"bookkeeping/config"
	"bookkeeping/dao"
	"bookkeeping/logic"
	"bookkeeping/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.MustGet(config.AuthMidUserNameKey).(string)
	sql_statement := `SELECT name, email FROM person WHERE id=$1`
	rows, err := dao.PostgresDB.Query(sql_statement, userID)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	var name string
	var email string
	for rows.Next() {
		if err := rows.Scan(&name, &email); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, nil)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"name":  name,
		"email": email,
	})

}

func AddPayment(c *gin.Context) {
	userID := c.MustGet(config.AuthMidUserNameKey).(string)
	var input model.PaymentRecord
	c.Bind(&input)
	if input.Class == "" || input.Date == "" || input.Payment == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	input.PersonID = userID
	if logic.CreatePayment(input) == false {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetPaymentHistory(c *gin.Context) {
	userID := c.MustGet(config.AuthMidUserNameKey).(string)
	tranlateMap := map[string]string{
		"eat":           "食",
		"clothes":       "衣",
		"live":          "住",
		"traffic":       "行",
		"educate":       "育",
		"entertainment": "樂",
	}
	sql_statement := `SELECT class, payment, date, remark FROM paymentrecord WHERE personid=$1`
	rows, err := dao.PostgresDB.Query(sql_statement, userID)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	retData := []model.PaymentRecord{}
	var tmpData model.PaymentRecord

	for rows.Next() {
		if err := rows.Scan(&tmpData.Class, &tmpData.Payment, &tmpData.Date, &tmpData.Remark); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, nil)
		}
		tmpData.Date = strings.ReplaceAll(tmpData.Date, "T00:00:00Z", "")
		tmpData.Class = tranlateMap[tmpData.Class]
		retData = append(retData, tmpData)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": retData,
	})
}

func GetAverage(c *gin.Context) {
	userID := c.MustGet(config.AuthMidUserNameKey).(string)
	currentYear := strconv.Itoa(time.Now().Year())
	currentMonth := time.Now().Month().String()
	transLateMonth := map[string]string{
		"January":   "1",
		"February":  "2",
		"March":     "3",
		"April":     "4",
		"May":       "5",
		"June":      "6",
		"July":      "7",
		"August":    "8",
		"September": "9",
		"October":   "10",
		"November":  "11",
		"December":  "12",
	}
	currentMonth = transLateMonth[currentMonth]
	// m-avg
	sql_statement := `SELECT SUM(payment)/COUNT(payment) AS monthAverage FROM paymentrecord WHERE personid=$1 AND DATE_PART('YEAR', date)=$2 GROUP BY DATE_PART('MONTH', date);`
	rows, err := dao.PostgresDB.Query(sql_statement, userID, currentYear)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	count := 0
	tmpNum := 0
	totalNum := 0
	for rows.Next() {
		if err := rows.Scan(&tmpNum); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		totalNum += tmpNum
		tmpNum = 0
		count++
	}
	monthAvg := strconv.Itoa(totalNum / count)
	// dayAvg
	sql_statement = `SELECT SUM(payment)/COUNT(payment) as dayAverage FROM paymentrecord WHERE personid=$1 AND DATE_PART('YEAR', date)=$2 AND DATE_PART('MONTH', date)=$3;`
	rows, err = dao.PostgresDB.Query(sql_statement, userID, currentYear, currentMonth)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	for rows.Next() {
		if err := rows.Scan(&tmpNum); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, nil)
			return
		}
	}
	dayAvg := strconv.Itoa(tmpNum)
	c.JSON(http.StatusOK, gin.H{
		"monthAvg": monthAvg,
		"dayAvg":   dayAvg,
	})
}
