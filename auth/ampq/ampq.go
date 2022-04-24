package ampq

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitAMPQ() (*amqp.Connection, func() error, error) {
	amqpServerURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	fmt.Printf("Connecting to AMPQ server at %s\n", amqpServerURL)
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)

	if err != nil {
		return nil, nil, err
	}

	return connectRabbitMQ, connectRabbitMQ.Close, nil
}
