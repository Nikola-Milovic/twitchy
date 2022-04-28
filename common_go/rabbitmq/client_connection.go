package rabbitmq

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second
)

//Maybe manage a pool of channels
type ClientConnection struct {
	logger        *zap.SugaredLogger
	connection    *amqp.Connection
	Channel       *amqp.Channel
	Done          chan os.Signal
	NotifyClose   chan *amqp.Error
	NotifyConfirm chan amqp.Confirmation
	IsConnected   bool
	Alive         bool
}

// NewClientConnection creates a new ClientConnection
func NewClientConnection(logger *zap.SugaredLogger, done chan os.Signal) *ClientConnection {
	return &ClientConnection{
		logger: logger,
		Done:   done,
		Alive:  true,
	}
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (c *ClientConnection) HandleReconnect(addr string, clientConnect func(*amqp.Channel) bool) {
	for c.Alive {
		c.IsConnected = false
		t := time.Now()
		fmt.Printf("Attempting to connect to rabbitMQ: %s\n", addr)
		var retryCount int
		for !c.connect(addr, clientConnect) {
			if !c.Alive {
				return
			}
			select {
			case <-c.Done:
				return
			case <-time.After(reconnectDelay + time.Duration(retryCount)*time.Second):
				c.logger.Info("disconnected from rabbitMQ and failed to connect")
				retryCount++
			}
		}
		c.logger.Infof("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
		select {
		case <-c.Done:
			return
		case <-c.NotifyClose:
		}
	}
}

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (c *ClientConnection) connect(addr string, clientConnect func(*amqp.Channel) bool) bool {
	conn, err := amqp.Dial(addr)
	if err != nil {
		c.logger.Errorf("failed to dial rabbitMQ server: %v", err)
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		c.logger.Errorf("failed connecting to channel: %v", err)
		return false
	}
	ch.Confirm(false)

	connectedClient := clientConnect(ch)

	if !connectedClient {
		c.logger.Errorf("failed to connect the client to rabbitmq: %v", err)
		return false
	}
	c.changeConnection(conn, ch)
	c.IsConnected = true
	return true
}

// changeConnection takes a new connection to the queue,
// and updates the channel listeners to reflect this.
func (c *ClientConnection) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.connection = connection
	c.Channel = channel
	c.NotifyClose = make(chan *amqp.Error)
	c.NotifyConfirm = make(chan amqp.Confirmation)
	c.Channel.NotifyClose(c.NotifyClose)
	c.Channel.NotifyPublish(c.NotifyConfirm)
}

func (c *ClientConnection) Close() error {
	err := c.Channel.Close()
	if err != nil {
		c.logger.Errorf("failed to close channel: %v", err)
	}
	err = c.connection.Close()
	if err != nil {
		return err
	}
	c.IsConnected = false

	return nil
}
