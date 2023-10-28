package commands

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"movie-downloader-bot/internal/client"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"regexp"
	"strings"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	cache   *redis.Client
	client  client.Client
	ctx     context.Context
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser, client client.Client) *Commander {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		client:  client,
		cache:   rdb,
		ctx:     context.Background(),
	}
}

func (cmd *Commander) HandleUpdate(update tgbotapi.Update) {
	//defer func() {
	//	if panicValue := recover(); panicValue != nil {
	//		log.Printf("Recovered from panic: %v", panicValue)
	//	}
	//}()

	// Handle callbacks
	if update.CallbackQuery != nil {

		clbk := strings.Split(update.CallbackQuery.Data, "|")
		cmdName := clbk[0]

		println(cmdName)

		switch cmdName {
		case "metamovie_download":
			log.Println("Download callback!")
			if len(clbk) != 2 {
				log.Println("No kpid!")
				return
			}
			cmd.DownloadByLinkOrId(update.CallbackQuery.Message, clbk[1], true)

		case "metamovie_torrents":
			log.Println("Torrent callback!")
			if len(clbk) != 2 {
				log.Fatal("No kpid!")
				return
			}
			cmd.SearchByLinkOrId(update.CallbackQuery.Message, clbk[1], true)

		case "movie_show":
			log.Println("Movie show callback!")
			if len(clbk) != 2 {
				log.Println("No cache id!")
				return
			}
			cmd.ShowMovie(update.CallbackQuery.Message, clbk[1])

		default:
			log.Println("Unknown callback!")
		}
		return
	}

	if update.Message == nil {
		return
	}

	// Handle static commands
	switch update.Message.Command() {
	case "start":
		cmd.Start(update.Message)
	}

	msgTxt := strings.ToLower(strings.Trim(update.Message.Text, " /"))
	downloadRe := regexp.MustCompile(params.Get().Commands.Download)
	searchRe := regexp.MustCompile(params.Get().Commands.Search)

	// Handle text commands
	switch {
	case strings.Contains(msgTxt, "kinopoisk.ru/film"):
		cmd.DownloadByLinkOrId(update.Message, msgTxt, false)

	case strings.Contains(msgTxt, "kinopoisk.ru/series"):
		cmd.SearchByLinkOrId(update.Message, msgTxt, false)

	case downloadRe.MatchString(msgTxt):
		cmd.SearchOrDownloadByTitle(update.Message, msgTxt, downloadRe, true)

	case searchRe.MatchString(msgTxt):
		cmd.SearchOrDownloadByTitle(update.Message, msgTxt, searchRe, false)
	}
}
