package main

import (
	"context"
	"math/rand"
	"net/http"
	"nikolamilovic/twitchy/accounts/ampq"
	"nikolamilovic/twitchy/auth/api"
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
	)
	rand.Seed(time.Now().UnixNano())

	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@172.28.0.1:5432/auth-db?sslmode=disable")
	os.Setenv("AMQP_SERVER_URL", "amqp://172.27.0.1:5672/")

	dbConn := db.InitDb()
	ampqConn := ampq.InitAMPQ()

	shutdowns = append(shutdowns, db.CloseDb(ctx, dbConn), ampqConn.Close)

	srv := api.NewServer(dbConn)

	server := http.Server{
		Addr:    ":4003",
		Handler: srv,
	}
	go gracefulShutdown(&server, shutdown, ctx)

	logger.Info("server starting: http://localhost" + server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}
}

func gracefulShutdown(server *http.Server, shutdown chan struct{}, ctx context.Context) {
	var (
		sigint = make(chan os.Signal, 1)
	)

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
