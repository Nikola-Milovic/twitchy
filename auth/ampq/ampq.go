package ampq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitAMPQ() (*amqp.Connection, func() error, error) {
	amqpServerURL := fmt.Sprintf("amqp://guest:guest@rabbitmq:5672/")

	fmt.Printf("Connecting to AMPQ server at %s\n", amqpServerURL)
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return nil, nil, err
	}

	return connectRabbitMQ, connectRabbitMQ.Close, nil
}
