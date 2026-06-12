package main

import (
	"fmt"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	defer fmt.Print("> ")
	return func(armyMv gamelogic.ArmyMove) {
		gs.HandleMove(armyMv)
	}
}
