package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	tgBotToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on %s", bot.Self.UserName)

	updateConfig := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		log.Printf("%+v\n", update)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You wrote: "+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
