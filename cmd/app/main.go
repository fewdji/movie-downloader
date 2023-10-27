package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"movie-downloader-bot/internal/commander"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"runtime"
	"strconv"
	"time"
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
	log.Printf("Authorized on %s", bot.Self.UserName)

	updateConfig := tgbotapi.UpdateConfig{
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(updateConfig)

	kpParser := meta.NewKpParser()
	tParser := torrent.NewJackettParser()
	commander := commands.NewCommander(bot, kpParser, tParser)

	//go MonitorTask()

	for update := range updates {
		commander.HandleUpdate(update)
	}
}

func MonitorTask() {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("task working")
		fmt.Println(runtime.NumGoroutine())
	}
}
