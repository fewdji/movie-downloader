package commands

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (cmd *Commander) DeleteMessage(chatID int64, messageId ...int) error {
	log.Println("Delete msg!")
	if messageId == nil {
		log.Println("No messageId!")
		return errors.New("no messageId")
	}
	for _, msgId := range messageId {
		_, err := cmd.bot.Request(tgbotapi.NewDeleteMessage(chatID, msgId))
		if err != nil {
			log.Println(err)
			//return err
		}
	}
	return nil
}
