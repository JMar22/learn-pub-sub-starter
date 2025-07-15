package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/JMar22/learn-pub-sub-starter/internal/pubsub"
	"github.com/JMar22/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	// Establish RabbitMQ connection
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	// Create RabbitMQ channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not open channel: %v", err)
	}
	defer ch.Close()
	fmt.Println("RabbitMQ channel established")

	// Publish pause message
	pauseMsg := routing.PlayingState{IsPaused: true}
	if err := pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		pauseMsg,
	); err != nil {
		log.Fatalf("could not publish pause message: %v", err)
	}
	fmt.Println("Published pause message")

	// Wait for interrupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down server...")
}
