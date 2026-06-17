package pubsub

import (
	"errors"

	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	DurableQueue   SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

type queueParams struct {
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

var queueTypes = map[SimpleQueueType]queueParams{
	DurableQueue: {
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       amqp.Table{"x-dead-letter-exchange": routing.DeadLetterExchange},
	},
	TransientQueue: {
		durable:    false,
		autoDelete: true,
		exclusive:  true,
		noWait:     false,
		args:       amqp.Table{"x-dead-letter-exchange": routing.DeadLetterExchange},
	},
}

func (s SimpleQueueType) queueParams() (queueParams, error) {
	params, exists := queueTypes[s]
	if !exists {
		return queueParams{}, errors.New("unknown queue type: " + string(s))
	}
	return params, nil
}

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
	Ignore
)
