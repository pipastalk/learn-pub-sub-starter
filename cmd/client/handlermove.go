package main

import (
	"fmt"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	defer fmt.Print("> ")
	return func(armyMv gamelogic.ArmyMove) pubsub.AckType {
		outcome := gs.HandleMove(armyMv)
		switch outcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, armyMv.Player.Username),
				gamelogic.RecognitionOfWar{
					Attacker: armyMv.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			return pubsub.NackDiscard
		}
	}
}
