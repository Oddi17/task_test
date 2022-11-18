package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	handlers "vito/pkg"
)

func main() {
	
	db, err := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()

	router.POST("/AddUserBalance", handlers.AddUserBalance)
	router.GET("/UserBalance/:id", handlers.GetUserBalance)
	router.PUT("/Transaction", handlers.Reserv)
	router.POST("/TransactionConfirm", handlers.TransactionConfirm)
	router.PUT("/TransactionReject", handlers.TransactionReject)

	router.Run()

}
