package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pipastalk/learn-pub-sub-starter/internal/gamelogic"
	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
	"github.com/pipastalk/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	var connectionStr = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionStr)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %s\n", err)
		os.Exit(1)
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ successfully!")
	ch, err := connection.Channel()
	if err != nil {
		fmt.Printf("Failed to create channel: %s\n", err)
		os.Exit(1)
	}
	defer ch.Close()
	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%s.*", routing.GameLogSlug),
		pubsub.SimpleQueueType("durable"),
	)
	if err != nil {
		fmt.Printf("Failed to declare and bind queue: %s\n", err)
		os.Exit(1)
	}
	err = pubsub.Subscribe(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%s.#", routing.GameLogSlug),
		pubsub.SimpleQueueType("durable"),
		HandlerWriteGameLog(),
		pubsub.HelperUnmarshallerGob[routing.GameLog](),
		case "pause":
			fmt.Println("Chill dude, We're pausing the game")
			err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Printf("Failed to publish JSON: %s\n", err)
				os.Exit(1)
			}
		case "resume":
			fmt.Println("Let's get back to it, resuming the game")
			err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Printf("Failed to publish JSON: %s\n", err)
				os.Exit(1)
			}
		case "quit":
			{
				fmt.Println("We all have to go sometime, goodbye")
				break
			}
		//endregion
		default:
			{
				fmt.Println("You seem to be chatting shit, (command unknown)")
			}
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Server has been terminated")

}
