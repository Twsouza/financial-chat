package services

import (
	"bot/application/dto"
	"bot/core"
	"encoding/json"
	"fmt"
	"log"
	"queue"
	"strings"
	"time"
)

type BotService struct {
	Command       map[string]func(string) (string, string, error)
	MessageBroker *queue.RabbitMQ
}

func NewBotService(messageBroker *queue.RabbitMQ) (*BotService, error) {
	return &BotService{
		Command:       initCommandMap(),
		MessageBroker: messageBroker,
	}, nil
}

func initCommandMap() map[string]func(string) (string, string, error) {
	commands := map[string]func(string) (string, string, error){}

	commands["stock"] = core.HandleStockCommand

	return commands
}

func (bs *BotService) HandleCommand(m *dto.Message) {
	prompt := strings.Split(m.Content, "=")
	if len(prompt) < 2 {
		bs.notify(m, "bot", "", fmt.Errorf("invalid command"))
		return
	}

	// the first element of prompt is the command
	cmd := strings.TrimPrefix(prompt[0], "/")
	// the rest of the prompt are the arguments
	args := strings.Join(prompt[1:], "=")

	fn, ok := bs.Command[cmd]
	if !ok {
		bs.notify(m, "bot", "", fmt.Errorf("command %s not found", cmd))
		return
	}

	botname, output, err := fn(args)
	bs.notify(m, botname, output, err)
}

func (b *BotService) notify(m *dto.Message, botname string, output string, err error) {
	newMessage := &dto.Message{
		Username:  botname,
		RoomID:    m.RoomID,
		CreatedAt: time.Now(),
	}

	if err != nil {
		newMessage.Content = fmt.Sprintf("Error: %s", err.Error())
	} else {
		newMessage.Content = output
	}

	notify, err := json.Marshal(newMessage)
	if err != nil {
		log.Printf("error marshalling message: %v", err)
	}

	err = b.MessageBroker.Notify(string(notify), "application/json", "amq.direct", "message")
	if err != nil {
		log.Printf("error notifying message: %v", err)
	}
}
