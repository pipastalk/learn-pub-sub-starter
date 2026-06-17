package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	messages, err := ch.Consume(
		queueName,
		"",
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}
	go func() {
		for msg := range messages {
			var payload T
			err = json.Unmarshal(msg.Body, &payload)
			switch handler(payload) {
			case Ack:
				fmt.Println("Acknowledged")
				msg.Ack(false)
			case NackRequeue:
				fmt.Println("Requeued")
				msg.Nack(false, true) // requeue the message
			case NackDiscard:
				fmt.Println("Discarded")
				msg.Nack(false, false) // discard the message
			case Ignore:
				// Do nothing, message will remain unacknowledged
			}
			// If ack is Ignore, do nothing (message will remain unacknowledged)
		}
	}()
	return nil
}
