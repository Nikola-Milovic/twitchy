package main

import (
	"context"
	"math/rand"
	"nikolamilovic/twitchy/accounts/api"
	"nikolamilovic/twitchy/accounts/service"
	db "nikolamilovic/twitchy/common/db"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
	shutdowns []func() error
)

func main() {
	var (
		shutdown = make(chan struct{})
		ctx      = context.Background()
		_        = make(chan os.Signal, 1)
	)

	rand.Seed(time.Now().UnixNano())

	dbConn, dbCleanup, err := db.InitDb(ctx, logger.Sugar().Named("db"))
	if err != nil {
		logger.Fatal("failed to init the db", zap.Error(err))
	}

	// amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
	// 	os.Getenv("RABBITMQ_USER"),
	// 	os.Getenv("RABBITMQ_PASSWORD"),
	// 	os.Getenv("RABBITMQ_HOST"),
	// 	os.Getenv("RABBITMQ_PORT"),
	// )

	// clientConnection := rabbitmq.NewClientConnection(logger.Sugar().Named("client_connection"), sigint)
	// client := client.New(amqpServerURL, logger.Sugar().Named("accounts_rabbitmq_client"), clientConnection)

	accountService := service.NewAccountService(dbConn)

	srv := api.NewServer(accountService)

	shutdowns = append(shutdowns, dbCleanup) //client.Close

	defer logger.Sync()

	go gracefulShutdown(srv.Server(), shutdown)

	logger.Debug("Go listening on 4002")

	err = srv.Listen(":4002")
	if err != nil {
		logger.Fatal("failed to start the server", zap.Error(err))
		os.Exit(1)
	}
}

func gracefulShutdown(server *fasthttp.Server, shutdown chan struct{}) {
	var (
		sigint = make(chan os.Signal, 1)
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}

	// close any other modules.
	for i := range shutdowns {
		shutdowns[i]()
	}

	close(shutdown)
}
