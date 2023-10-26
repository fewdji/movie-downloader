package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (cmd *Commander) Start(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.StartMsgTxt)
	msg.ParseMode = "markdown"
	_, err := cmd.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
