package main

import (
	"context"
	"math/rand"
	"nikolamilovic/twitchy/accounts/ampq"
	"nikolamilovic/twitchy/accounts/api"
	"nikolamilovic/twitchy/accounts/db"
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
	)

	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@172.28.0.1:5432/accounts-db?sslmode=disable")
	os.Setenv("AMQP_SERVER_URL", "amqp://172.27.0.1:5672/")

	rand.Seed(time.Now().UnixNano())

	dbConn, dbCleanup, err := db.InitDb(ctx)
	_, ampqCleanup, err := ampq.InitAMPQ()

	shutdowns = append(shutdowns, dbCleanup, ampqCleanup)

	if err != nil {
		logger.Fatal("Unable to initialize the app", zap.Error(err))
		os.Exit(1)
	}

	srv := api.NewServer(dbConn)

	go gracefulShutdown(srv.Server(), shutdown)

	logger.Debug("Go listening on 4002")

	srv.Listen(":4002")

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
