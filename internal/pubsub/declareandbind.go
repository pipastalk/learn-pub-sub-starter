package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to create channel: %w", err)
	}
	qParams, err := queueType.queueParams()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("invalid queue type: %w", err)
	}
	queue, err := ch.QueueDeclare(
		queueName,
		qParams.durable,
		qParams.autoDelete,
		qParams.exclusive,
		qParams.noWait,
		qParams.args,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare queue: %w", err)
	}
	err = ch.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to bind queue: %w", err)
	}
	return ch, queue, nil
}
