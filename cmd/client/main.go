package main

import (
	"fmt"
	"os"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	var connectionStr = "amqp://guest:guest@localhost:5672/"
	fmt.Printf("Connecting to RabbitMQ at %s...\n", connectionStr)
	connection, err := amqp.Dial(connectionStr)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		os.Exit(1)
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ successfully!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("Unable to get client username: %s\n", err)
		os.Exit(1)
	}
	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueType("transient"),
	)
	if err != nil {
		fmt.Printf("Failed to declare and bind queue: %s\n", err)
		os.Exit(1)
	}
	gameState := gamelogic.NewGameState(username)
replLoop:
	for {
		words := gamelogic.GetInput()
		//commands
		switch words[0] {
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Error occurred while spawning: %s\n", err)
			}

		case "move":
			mv, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Printf("Error occurred while moving: %s\n", err)
			}
			unitsString := ""
			if len(mv.Units) > 1 {
				unitsString = fmt.Sprintf("%d units", len(mv.Units))
			} else {
				unitsString = fmt.Sprintf("%s", mv.Units[0].Rank)
			}
			fmt.Printf("Moving %s to %s...\n", unitsString, mv.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break replLoop
		default:
			fmt.Printf("That random collection of symbols is noncesense, (Unknown command): %s\n", words[0])

		}
	}
	fmt.Printf("Client %s has been terminated\n", username)
}
