package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"movie-downloader-bot/internal/cache/redis"
	"movie-downloader-bot/internal/client/qbittorrent"
	"movie-downloader-bot/internal/commander"
	"movie-downloader-bot/internal/parser/meta/kpunofficial"
	"movie-downloader-bot/internal/parser/torrent/jackett"
	"movie-downloader-bot/internal/storage/postgres"
	tracking "movie-downloader-bot/internal/tracker"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug, _ = strconv.ParseBool(os.Getenv("TG_BOT_DEBUG"))
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})

	rCache := redis.NewRedis()
	kpParser := kpunofficial.NewKpParser(rCache)
	tParser := jackett.NewJackettParser()
	qClient := qbittorrent.NewQbittorrent()
	pStorage := postgres.NewPostgres()

	tracker := tracking.NewTracker(
		kpParser,
		tParser,
		qClient,
		pStorage,
	)

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
