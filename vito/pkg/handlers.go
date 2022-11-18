package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// создание или пополнение баланса
func AddUserBalance(c *gin.Context) {
	var User Balances
	err := c.Bind(&User)
	if err != nil {
		log.Fatalln(err)
	}
	account, err := User.createBalance()
	if err != nil {
		if fmt.Sprintf("%s", err) == "bad payment" {
			message := fmt.Sprintf("Error : %s", err) 
			c.JSON(http.StatusPaymentRequired, gin.H{
				"message": message,
			})
			return
		} else {
			log.Fatalln(err)
		}
	}
	message := fmt.Sprintf(
		"successfully completed! user id:%d ; balance:%0.2f ",
		User.ID, account,
	)
	//c.IndentedJSON()
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
	})
}

// Получение баланса пользователя
func GetUserBalance(c *gin.Context) {
	//var b Balances
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalln(err)
	}
	user, err := UserBalance(id)
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, user)
}

// Резервирование
func Reserv(c *gin.Context) {
	var t Reservation
	err := c.Bind(&t)
	if err != nil {
		log.Fatalln(err)
	}
	newvalue, err := t.TransactionReserv()
	if err != nil {
		if fmt.Sprintf("%s", err) == "insufficient funds" {
			message := fmt.Sprintf("Error : %s", err) 
			c.JSON(http.StatusPaymentRequired, gin.H{
				"message": message,
			})
			return
		} else {
			log.Fatalln(err)
		}
	}
	message := fmt.Sprintf("Successfully reserved! id:%d current balance:%0.2f ", t.ID, newvalue)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
	return
}

// Подтверждение транзакции
func TransactionConfirm(c *gin.Context) {
	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		log.Fatalln(err)
	}
	err = t.Confirm() 
	if err != nil {
		log.Fatalln(err)
	}
	message := fmt.Sprintf("Transaction of order:%d was confirmed! ", t.ID_ORDER)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})

}

// Отклонение транзакции
func TransactionReject(c *gin.Context) {
	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		log.Fatalln(err)
	}
	err = t.Reject() 
	if err != nil {
		log.Fatalln(err)
	}
	message := fmt.Sprintf("Transaction of order:%d was rejected! ", t.ID_ORDER)
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})

}
