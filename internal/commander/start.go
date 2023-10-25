package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func (c *Commander) Start(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, os.Getenv("START_MSG"))
	msg.ParseMode = "markdown"
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
