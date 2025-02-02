package main

import (
	"context"
	"finance-api/config"
	"finance-api/internal/handler"
	"finance-api/internal/repository"
	"finance-api/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close(context.Background())

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()

	r.POST("/deposit", h.Deposit)
	r.POST("/transfer", h.Transfer)
	r.GET("/transactions", h.GetTransactions)

	r.Run(":8080")
}
