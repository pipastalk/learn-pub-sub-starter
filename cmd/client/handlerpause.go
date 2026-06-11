package main

import (
	"fmt"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	defer fmt.Print("> ")
	return func(ps routing.PlayingState) {
		gs.HandlePause(ps)
	}
}
