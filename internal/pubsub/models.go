package pubsub

import (
	"errors"

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
	DurableQueue: queueParams{
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	},
	TransientQueue: queueParams{
		durable:    false,
		autoDelete: true,
		exclusive:  true,
		noWait:     false,
		args:       nil,
	},
}

func (s SimpleQueueType) queueParams() (queueParams, error) {
	params, exists := queueTypes[s]
	if !exists {
		return queueParams{}, errors.New("unknown queue type: " + string(s))
	}
	return params, nil
}
