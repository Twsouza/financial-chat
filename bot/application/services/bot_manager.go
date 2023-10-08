package services

import (
	"bot/application/dto"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
)

type BotManager struct {
	Message    chan amqp.Delivery
	BotService *BotService
}

func NewBotManager(message chan amqp.Delivery, bs *BotService) (*BotManager, error) {
	return &BotManager{
		Message:    message,
		BotService: bs,
	}, nil
}

func (b *BotManager) Run() {
	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading CONCURRENCY_WORKERS env var: %s", os.Getenv("CONCURRENCY_WORKERS"))
	}

	wg := sync.WaitGroup{}
	for qty := 0; qty < concurrency; qty++ {
		wg.Add(1)
		go b.Worker()
	}

	fmt.Println("Bot manager running")
	wg.Wait()

	close(b.Message)
}

func (bm *BotManager) Worker() {
	for message := range bm.Message {
		commandMessage := &dto.Message{}
		err := json.Unmarshal(message.Body, commandMessage)
		if err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		bm.BotService.HandleCommand(commandMessage)
		message.Ack(false)
	}
}
