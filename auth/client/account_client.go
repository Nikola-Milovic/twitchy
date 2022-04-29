package client

import (
	"encoding/json"
	"errors"
	"fmt"
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
	PublishAccountCreatedEvent(event event.AccountCreatedEventData) error
}

// AccountClient holds necessery information for rabbitMQ
type AccountClient struct {
	logger     *zap.SugaredLogger
	connection *rabbitmq.ClientConnection
	threads    int
	wg         *sync.WaitGroup
}

func New(addr string, l *zap.SugaredLogger, connection *rabbitmq.ClientConnection) *AccountClient {
	threads := runtime.GOMAXPROCS(0)
	if numCPU := runtime.NumCPU(); numCPU > threads {
		threads = numCPU
	}

	client := AccountClient{
		logger:     l,
		threads:    threads,
		connection: connection,
		wg:         &sync.WaitGroup{},
	}

	go client.connection.HandleReconnect(addr, client.connect)

	return &client
}

// Push a new message that an account has been created
func (c *AccountClient) PublishAccountCreatedEvent(data event.AccountCreatedEventData) error {
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	baseEv := event.BaseEvent{
		Type:    event.AccountCreatedType,
		Payload: string(payload),
	}

	ev, err := json.Marshal(baseEv)

	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	return c.push(constants.AccountCreatedKey, ev)
}

// Push will push data onto the queue, and wait for a confirmation.
// If no confirms are received until within the resendTimeout,
// it continuously resends messages until a confirmation is received.
// This will block until the server sends a confirm

//TODO add a timeout to the push and store the event into db, this shouldn't block
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

// unsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
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

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (c *AccountClient) connect(ch *amqp.Channel) bool {

	err := ch.ExchangeDeclare(constants.AccountsExchange, "topic", true, false, false, false, nil)

	if err != nil {
		c.logger.Errorf("failed to declare exchange: %v", err)
		return false
	}

	_, err = ch.QueueDeclare(
		constants.AuthServiceQueue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		c.logger.Errorf("failed to declare %s queue: %v", constants.AuthServiceQueue, err)
		return false
	}

	_, err = ch.QueueDeclare(
		constants.AccountsQueue,
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
	err = ch.QueueBind(constants.AuthServiceQueue, "", constants.AccountsExchange, true, nil)

	if err != nil {
		c.logger.Errorf("failed to bind stream queue: %v", err)
		return false
	}

	return true
}

// func (c *AccountClient) stream(cancelCtx context.Context) error {
// 	c.wg.Add(c.threads)

// 	for {
// 		if c.isConnected {
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}

// 	err := c.channel.Qos(1, 0, false)
// 	if err != nil {
// 		return err
// 	}

// 	var connectionDropped bool

// 	for i := 1; i <= c.threads; i++ {
// 		msgs, err := c.channel.Consume(
// 			c.streamQueue,
// 			consumerName(i), // Consumer
// 			false,           // Auto-Ack
// 			false,           // Exclusive
// 			false,           // No-local
// 			false,           // No-Wait
// 			nil,             // Args
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		go func() {
// 			defer c.wg.Done()
// 			for {
// 				select {
// 				case <-cancelCtx.Done():
// 					return
// 				case msg, ok := <-msgs:
// 					if !ok {
// 						connectionDropped = true
// 						return
// 					}
// 					c.parseEvent(msg)
// 				}
// 			}
// 		}()

// 	}

// 	c.wg.Wait()

// 	if connectionDropped {
// 		return ErrDisconnected
// 	}

// 	return nil
// }

// type event struct {
// 	Job  string `json:"job"`
// 	Data string `json:"data"`
// }

// func (c *AccountClient) parseEvent(msg amqp.Delivery) {
// 	l := c.logger.Named("parseEvent")
// 	startTime := time.Now()

// 	var evt event
// 	err := json.Unmarshal(msg.Body, &evt)
// 	if err != nil {
// 		logAndNack(msg, l, startTime, "unmarshalling body: %s - %s", string(msg.Body), err.Error())
// 		return
// 	}

// 	if evt.Data == "" {
// 		logAndNack(msg, l, startTime, "received event without data")
// 		return
// 	}

// 	defer func(e event, m amqp.Delivery, logger *zap.SugaredLogger) {
// 		if err := recover(); err != nil {
// 			stack := make([]byte, 8096)
// 			stack = stack[:runtime.Stack(stack, false)]
// 			logger.Error("panic recovery for rabbitMQ message")
// 			msg.Nack(false, false)
// 		}
// 	}(evt, msg, l)

// 	// switch evt.Job {
// 	// case "job1":
// 	//     // Call an actual function
// 	//     err = func()
// 	// case "job1":
// 	//     err = func()
// 	// default:
// 	//     msg.Reject(false)
// 	//     return
// 	// }

// 	if err != nil {
// 		logAndNack(msg, l, startTime, err.Error())
// 		return
// 	}

// 	l.Errorf("Took ms %d, succeeded %s", time.Since(startTime).Milliseconds(), evt.Job)
// 	msg.Ack(false)
// }

// func logAndNack(msg amqp.Delivery, l *zap.SugaredLogger, t time.Time, err string, args ...interface{}) {
// 	msg.Nack(false, false)
// 	l.Errorf("Took ms %d, %e", time.Since(t).Milliseconds(), err)
// }

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
