package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"regexp"
	"strings"
	"time"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	cache   *cache.Cache
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser) *Commander {
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		cache:   cache.New(15*time.Minute, 30*time.Minute),
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
