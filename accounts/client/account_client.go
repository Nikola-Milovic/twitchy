package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nikolamilovic/twitchy/accounts/service"
	"nikolamilovic/twitchy/common/constants"
	"nikolamilovic/twitchy/common/event"
	"nikolamilovic/twitchy/common/rabbitmq"
	"runtime"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

//https://www.ribice.ba/golang-rabbitmq-client/
type IAccountClient interface {
}

// AccountClient holds necessery information for rabbitMQ
type AccountClient struct {
	service    service.IAccountService
	logger     *zap.SugaredLogger
	connection *rabbitmq.ClientConnection
	threads    int
	wg         *sync.WaitGroup
}

func New(addr string, l *zap.SugaredLogger, service service.IAccountService, connection *rabbitmq.ClientConnection) *AccountClient {
	threads := runtime.GOMAXPROCS(0)
	if numCPU := runtime.NumCPU(); numCPU > threads {
		threads = numCPU
	}

	client := AccountClient{
		logger:     l,
		service:    service,
		threads:    threads,
		connection: connection,
		wg:         &sync.WaitGroup{},
	}

	go client.connection.HandleReconnect(addr, client.connect)
	return &client
}

func (c *AccountClient) push(key string, data []byte) error {
	if !c.connection.IsConnected {
		return errors.New("failed to push push: not connected")
	}
	for {
		err := c.unsafePush(key, data)
		if err != nil {
			if err == rabbitmq.ErrDisconnected {
				continue
			}
			return err
		}
		select {
		case confirm := <-c.connection.NotifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(resendDelay):
		}
	}
}

func (c *AccountClient) unsafePush(key string, data []byte) error {
	if !c.connection.IsConnected {
		return rabbitmq.ErrDisconnected
	}

	return c.connection.Channel.Publish(
		constants.AccountsExchange, // Exchange
		key,                        // Routing key
		false,                      // Mandatory
		false,                      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

func (c *AccountClient) connect(ch *amqp.Channel) bool {

	err := ch.ExchangeDeclare(constants.AccountsExchange, "topic", true, false, false, false, nil)

	if err != nil {
		c.logger.Errorf("failed to declare exchange: %v", err)
		return false
	}

	_, err = ch.QueueDeclare(
		constants.AccountsQueue+"_accounts",
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		c.logger.Errorf("failed to declare %s queue: %v", constants.AccountsQueue, err)
		return false
	}

	err = ch.QueueBind(constants.AccountsQueue, constants.AccountCreatedKey, constants.AccountsExchange, true, nil)
	if err != nil {
		c.logger.Errorf("failed to bind push queue: %v", err)
		return false
	}

	return true
}

func (c *AccountClient) stream(cancelCtx context.Context) error {
	c.wg.Add(c.threads)

	for {
		if c.connection.IsConnected {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err := c.connection.Channel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	var connectionDropped bool

	for i := 1; i <= c.threads; i++ {
		msgs, err := c.connection.Channel.Consume(
			constants.AccountsQueue,
			consumerName(i), // Consumer
			false,           // Auto-Ack
			false,           // Exclusive
			false,           // No-local
			false,           // No-Wait
			nil,             // Args
		)
		if err != nil {
			return err
		}

		go func() {
			defer c.wg.Done()
			for {
				select {
				case <-cancelCtx.Done():
					return
				case msg, ok := <-msgs:
					if !ok {
						connectionDropped = true
						return
					}
					c.parseEvent(msg)
				}
			}
		}()

	}

	c.wg.Wait()

	if connectionDropped {
		return rabbitmq.ErrDisconnected
	}

	return nil
}

func (c *AccountClient) parseEvent(msg amqp.Delivery) {
	l := c.logger.Named("parseEvent")
	startTime := time.Now()

	var evt event.BaseEvent
	err := json.Unmarshal(msg.Body, &evt)
	if err != nil {
		logAndNack(msg, l, startTime, "unmarshalling body: %s - %s", string(msg.Body), err.Error())
		return
	}

	if evt.Payload == "" {
		logAndNack(msg, l, startTime, "received event without data")
		return
	}

	defer func(e event.BaseEvent, m amqp.Delivery, logger *zap.SugaredLogger) {
		if err := recover(); err != nil {
			stack := make([]byte, 8096)
			stack = stack[:runtime.Stack(stack, false)]
			logger.Error("panic recovery for rabbitMQ message")
			msg.Nack(false, false)
		}
	}(evt, msg, l)

	switch evt.Type {
	case event.AccountCreatedType:
		payload := &event.AccountCreatedEventData{}
		err := json.Unmarshal([]byte(evt.Payload), payload)
		if err != nil {
			logAndNack(msg, l, startTime, "failed to parse event data")
		}
		c.service.CreateUser(*payload)
	default:
		msg.Reject(false)
		return
	}

	if err != nil {
		logAndNack(msg, l, startTime, err.Error())
		return
	}

	l.Infof("Took ms %d, succeeded %s", time.Since(startTime).Milliseconds(), evt.Type)
	msg.Ack(false)
}

func logAndNack(msg amqp.Delivery, l *zap.SugaredLogger, t time.Time, err string, args ...interface{}) {
	msg.Nack(false, false)
	l.Errorf("Took ms %d, %e", time.Since(t).Milliseconds(), err)
}

func (c *AccountClient) Close() error {
	if !c.connection.IsConnected {
		return nil
	}
	c.connection.Alive = false
	c.logger.Info("Waiting for current messages to be processed...")
	c.wg.Wait()
	for i := 1; i <= c.threads; i++ {
		fmt.Println("Closing consumer: ", i)
		err := c.connection.Channel.Cancel(consumerName(i), false)
		if err != nil {
			return fmt.Errorf("error canceling consumer %s: %v", consumerName(i), err)
		}
	}

	err := c.connection.Close()

	if err != nil {
		return err
	}

	c.logger.Info("gracefully stopped rabbitMQ connection")
	return nil
}

func consumerName(i int) string {
	return fmt.Sprintf("go-consumer-%v", i)
}
