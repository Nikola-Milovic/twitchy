package ampq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitAMPQ() (*amqp.Connection, func() error, error) {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return nil, nil, err
	}

	return connectRabbitMQ, connectRabbitMQ.Close, nil
}
