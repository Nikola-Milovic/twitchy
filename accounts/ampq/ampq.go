package ampq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitAMPQ() *amqp.Connection {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	return connectRabbitMQ
}
