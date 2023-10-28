package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (cmd *Commander) Start(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.StartMsgTxt)
	msg.ParseMode = "markdown"
	cmd.bot.Send(msg)
}
