package main

import (
	"fmt"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	defer fmt.Print("> ")
	return func(ps routing.PlayingState) pubsub.AckType {
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}
