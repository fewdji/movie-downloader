package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	params "movie-downloader-bot/internal/config"
)

func (cmd *Commander) Start(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.StartMsgTxt)
	msg.ParseMode = "markdown"
	cmd.bot.Send(msg)
}
