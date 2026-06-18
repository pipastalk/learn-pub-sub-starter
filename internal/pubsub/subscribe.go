package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func HelperUnmarshallerJSON() func(body []byte) (any, error) {
	return func(body []byte) (any, error) {
		var payload any
		err := json.Unmarshal(body, &payload)
		if err != nil {
			var zero any
			return zero, err
		}
		return payload, nil
	}
}

func HelperUnmarshallerGob() func([]byte) (any, error) {
	return func(body []byte) (any, error) {
		var buf bytes.Buffer
		buf.Write(body)
		dec := gob.NewDecoder(&buf)
		var payload any
		err := dec.Decode(&payload)
		if err != nil {
			var zero any
			return zero, err
		}
		return payload, nil
	}
}

func Subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
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
			payload, _ := unmarshaller(msg.Body) //TODO create error handling
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
