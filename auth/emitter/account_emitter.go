package emitter

import (
	"encoding/json"
	"fmt"
	"log"
	"nikolamilovic/twitchy/auth/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IAccountEmitter interface {
	Emit(event model.AccountCreatedEvent) error
}

// AccountEmitter for publishing AMQP events
type AccountEmitter struct {
	connection *amqp.Connection
}

var (
	accountExchangeName = "account_exchange"
)

func (e *AccountEmitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()
	return accountExchange(channel)
}

// Push (Publish) a specified message to the AMQP exchange
func (e *AccountEmitter) Emit(event model.AccountCreatedEvent) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	data, err := json.Marshal(event)

	if err != nil {
		return err
	}

	err = channel.Publish(
		accountExchangeName,
		"account.created",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	log.Printf("Sending message: %v -> %s", event, accountExchangeName)
	return nil
}

func NewAccountEmitter(conn *amqp.Connection) (IAccountEmitter, error) {
	emitter := AccountEmitter{
		connection: conn,
	}

	if conn.IsClosed() {
		return nil, fmt.Errorf("connection is closed")
	}

	err := emitter.setup()
	if err != nil {
		return nil, err
	}

	return &emitter, nil
}

func accountCreatedQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"account.created", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
}

func accountExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		accountExchangeName, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
}
