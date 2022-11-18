package main

import (
	"BalanceService/internal/balance"
	"BalanceService/internal/balance/db"
	"BalanceService/internal/config"
	"BalanceService/pkg/client/postgresql"
	"BalanceService/pkg/logging"
	"context"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("123")
	server := gin.Default()

	cfg := config.GetConfig()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	database := db.NewDB(client, logger)

	handler := balance.NewHandler(database, logger)
	handler.Register(server)
	server.Run(":8999")
}
