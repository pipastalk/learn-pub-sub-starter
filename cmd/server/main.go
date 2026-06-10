package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pipastalk/learn-pub-sub-starter/internal/pubsub"
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
	err = pubsub.PublishJSON(ch, ExchangePerilDirect, PauseKey, PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("Failed to publish JSON: %s\n", err)
		os.Exit(1)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Server has been terminated")

}
