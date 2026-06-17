package main

import (
	"fmt"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
	"github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)
		logMsg := ""
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeOpponentWon:
			logMsg = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeYouWon:
			logMsg = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeDraw:
			logMsg = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
		default:
			return pubsub.NackDiscard
		}
		err := pubsub.PublishGameLog(ch, logMsg, rw.Attacker, gs.Player)
		if err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
