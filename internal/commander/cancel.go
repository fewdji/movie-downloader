package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (cmd *Commander) DeleteMessage(chatID int64, messageId ...int) {
	if messageId == nil {
		return
	}
	for _, msgId := range messageId {
		_, err := cmd.bot.Request(tgbotapi.NewDeleteMessage(chatID, msgId))
		if err != nil {
			log.Println("DeleteMessage:", err)
		}
	}
}
