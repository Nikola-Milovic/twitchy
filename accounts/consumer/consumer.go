package consumer

import (
	"nikolamilovic/twitchy/accounts/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AccountConsumer struct {
	conn *amqp.Connection
}

func NewAccountConsumer(conn *amqp.Connection) *AccountConsumer {
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	ch.ExchangeDeclare("account_topic", "topic", true, false, false, false, nil)

	// q, err := ch.QueueDeclare(
	// 	"account.created", // name
	// 	true,              // durable
	// 	false,             // delete when unused
	// 	false,             // exclusive
	// 	false,             // no-wait
	// 	nil,               // arguments
	// )
	// utils.FailOnError(err, "Failed to declare a queue")

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	false,  // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// utils.FailOnError(err, "Failed to register a consumer")

	return &AccountConsumer{conn: conn}
}
