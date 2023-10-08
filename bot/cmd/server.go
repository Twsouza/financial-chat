package main

import (
	"bot/application/services"
	"log"
	"queue"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()
	rabbitMQ.Consume(messageChannel)

	bs, err := services.NewBotService(rabbitMQ)
	if err != nil {
		log.Fatalf("error creating bot service: %v", err)
	}

	botManager, err := services.NewBotManager(messageChannel, bs)
	if err != nil {
		log.Fatalf("error creating bot: %v", err)
	}

	botManager.Run()
}
