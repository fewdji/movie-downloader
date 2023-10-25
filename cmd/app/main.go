package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"movie-downloader-bot/internal/commander"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load(".env", "config.env")
	if err != nil {
		log.Fatal(err)
	}

	pm := params.NewParams()

	kpParser := meta.NewKpParser()
	tp := torrent.NewJackettParser()

	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug, err = strconv.ParseBool(os.Getenv("TG_BOT_DEBUG"))
	log.Printf("Authorized on %s", bot.Self.UserName)

	updateConfig := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(updateConfig)

	commander := commands.NewCommander(bot, kpParser, tp, *pm)

	for update := range updates {
		log.Printf("%+v\n", update)
		commander.HandleUpdate(update)
	}
}
