package pubsub

import (
	"fmt"
	"time"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(ch *amqp091.Channel, logMsg string, attacker, owner gamelogic.Player) error {
	gl := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     logMsg,
		Username:    owner.Username,
	}
	err := PublishGob(
		ch,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, attacker.Username),
		gl,
	)
	if err != nil {
		return err
	}
	return nil
}
