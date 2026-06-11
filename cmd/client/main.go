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
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("Client %s has been terminated\n", username)
}
