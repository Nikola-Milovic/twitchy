package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"nikolamilovic/twitchy/auth/api"
	"nikolamilovic/twitchy/auth/client"
	"nikolamilovic/twitchy/auth/db"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		sigint   = make(chan os.Signal, 1)
	)
	rand.Seed(time.Now().UnixNano())

	dbConn, dbCleanup, err := db.InitDb(ctx, logger.Sugar().Named("db"))
	if err != nil {
		logger.Fatal("failed to init the db", zap.Error(err))
	}
	// ampq, ampqCleanup, err := ampq.InitAMPQ(logger.Sugar().Named("ampq"))
	// if err != nil {
	// 	logger.Fatal("failed to connect to", zap.Error(err))
	// }
	amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	client := client.New("", "accounts.created", amqpServerURL, logger.Sugar().Named("rabbitmq"), sigint)

	srv, err := api.NewServer(dbConn, client)

	if err != nil {
		logger.Fatal("Unable to initialize the server", zap.Error(err))
		os.Exit(1)
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	server := http.Server{
		Addr:    port,
		Handler: srv,
	}

	shutdowns = append(shutdowns, dbCleanup)

	defer logger.Sync()

	go gracefulShutdown(&server, shutdown, ctx, sigint)

	logger.Info("Server starting and listening at port " + server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
}

func gracefulShutdown(server *http.Server, shutdown chan struct{}, ctx context.Context, sigint chan os.Signal) {
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}

	// close any other modules.
	for i := range shutdowns {
		shutdowns[i]()
	}

	close(shutdown)
}
