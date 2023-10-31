package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"movie-downloader-bot/internal/cache"
	"movie-downloader-bot/internal/client"
	"movie-downloader-bot/internal/commander"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"strconv"
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

	bot.Debug, err = strconv.ParseBool(os.Getenv("TG_BOT_DEBUG"))

	updateConfig := tgbotapi.UpdateConfig{}
	updates := bot.GetUpdatesChan(updateConfig)

	cache := cache.NewRedis()
	kpParser := meta.NewKpParser(cache)
	tParser := torrent.NewJackettParser()
	client := client.NewQbittorrent()

	commander := commands.NewCommander(
		bot,
		kpParser,
		tParser,
		client,
		cache,
	)

	//tasker := tasks.NewTasker()
	//go tasker.Monitor()

	for update := range updates {
		commander.HandleUpdate(update)
	}
}
