package pubsub

import (
	"errors"
)

type SimpleQueueType string

const (
	DurableQueue   SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

// when using queueParams be sure to add args manually after
type queueParams struct {
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
}

var queueTypes = map[SimpleQueueType]queueParams{
	DurableQueue: {
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
	},
	TransientQueue: {
		durable:    false,
		autoDelete: true,
		exclusive:  true,
		noWait:     false,
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
