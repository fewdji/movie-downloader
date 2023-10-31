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
	"movie-downloader-bot/internal/storage"
	"movie-downloader-bot/internal/tracker"
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

	rCache := cache.NewRedis()
	kpParser := meta.NewKpParser(rCache)
	tParser := torrent.NewJackettParser()
	qClient := client.NewQbittorrent()
	storage := storage.NewPostgres()
	tracker := tracker.NewTracker(kpParser, tParser, qClient, storage)

	commander := commands.NewCommander(
		bot,
		kpParser,
		tParser,
		qClient,
		tracker,
		rCache,
	)

	go tracker.Run()

	for update := range updates {
		commander.HandleUpdate(update)
	}
}
