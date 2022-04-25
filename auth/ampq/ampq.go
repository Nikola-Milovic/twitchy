package ampq

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func InitAMPQ(logger *zap.SugaredLogger) (*amqp.Connection, func() error, error) {
	amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)

	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		if err == nil {
			break
		} else {
			conn, err = amqp.Dial(amqpServerURL)
		}

		logger.Error("Rabbitmq connection error, %s\n", err.Error())
		logger.Info("Retrying connection to rabbitmq...\n	defer logger.Sync()")
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if err != nil {
		return nil, nil, err
	}

	logger.Info("Connected to AMPQ server at %s\n", amqpServerURL)

	return conn, conn.Close, nil
}
