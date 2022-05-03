package client

import (
	"nikolamilovic/twitchy/accounts/service/mock"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func TestParseEventAck(t *testing.T) {
	//Set up test
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client := &AccountClient{
		logger:  zap.L().Sugar().Named("test"),
		service: &mock.AccountServiceMock{},
	}

	ack := NewMockAcknowledger(ctl)

	ack.EXPECT().Ack(gomock.Any(), false)

	//WHEN
	client.parseEvent(
		amqp091.Delivery{
			Acknowledger: ack,
			ContentType:  "application/json",
			Body: []byte(`{
 	  "type":"account_created",
 	  "payload":"{\"id\":12345,\"email\":\"test@gmail.com\"}"
		}`),
		},
	)
}
func TestParseEventNackWhenNoPayload(t *testing.T) {
	//Set up test
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client := &AccountClient{
		logger:  zap.L().Sugar().Named("test"),
		service: &mock.AccountServiceMock{},
	}

	ack := NewMockAcknowledger(ctl)

	ack.EXPECT().Nack(gomock.Any(), false, false)

	//WHEN
	client.parseEvent(
		amqp091.Delivery{
			Acknowledger: ack,
			ContentType:  "application/json",
			Body: []byte(`{
 	  "type":"account_created"
	   }`),
		},
	)
}
