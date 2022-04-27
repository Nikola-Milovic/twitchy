package client

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

//Maybe manage a pool of channels
type ClientConnection struct {
	logger        *zap.SugaredLogger
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan os.Signal
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	isConnected   bool
	alive         bool
}

// NewClientConnection creates a new ClientConnection
func NewClientConnection(logger *zap.SugaredLogger, done chan os.Signal) *ClientConnection {

	return &ClientConnection{
		logger: logger,
		done:   done,
		alive:  true,
	}
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (c *ClientConnection) handleReconnect(addr string, clientConnect func(*amqp.Channel) bool) {
	for c.alive {
		c.isConnected = false
		t := time.Now()
		fmt.Printf("Attempting to connect to rabbitMQ: %s\n", addr)
		var retryCount int
		for !c.connect(addr, clientConnect) {
			if !c.alive {
				return
			}
			select {
			case <-c.done:
				return
			case <-time.After(reconnectDelay + time.Duration(retryCount)*time.Second):
				c.logger.Info("disconnected from rabbitMQ and failed to connect")
				retryCount++
			}
		}
		c.logger.Infof("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
		select {
		case <-c.done:
			return
		case <-c.notifyClose:
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
	c.isConnected = true
	return true
}

// changeConnection takes a new connection to the queue,
// and updates the channel listeners to reflect this.
func (c *ClientConnection) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.connection = connection
	c.channel = channel
	c.notifyClose = make(chan *amqp.Error)
	c.notifyConfirm = make(chan amqp.Confirmation)
	c.channel.NotifyClose(c.notifyClose)
	c.channel.NotifyPublish(c.notifyConfirm)
}
