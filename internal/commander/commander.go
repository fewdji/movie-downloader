package commands

import (
	"context"
	"encoding/json"
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

type CommandData struct {
	MessageId int    `json:"m"`
	Command   string `json:"c"`
	Key       string `json:"k"`
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

		cmdData := CommandData{}
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &cmdData)
		if err != nil {
			log.Println(err)
			return
		}

		switch cmdData.Command {
		case "mm_down":
			cmd.DownloadBest(update.CallbackQuery.Message, cmdData)

		case "mm_tor":
			log.Println("Torrent callback!")
			cmd.ShowMovieList(update.CallbackQuery.Message, cmdData)

		case "m_sh":
			log.Println("Movie show callback!")
			cmd.ShowMovie(update.CallbackQuery.Message, cmdData)

		case "dl_f", "dl_s", "dl_t", "dl_w":
			cmd.DownloadMovie(update.CallbackQuery.Message, cmdData)

		default:
			log.Println("Unknown callback:", cmdData.Command)
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

	msgTxt := strings.ToLower(update.Message.Text)
	cmdData := CommandData{Key: msgTxt}
	downloadRe := regexp.MustCompile(params.Get().Commands.Download)
	searchRe := regexp.MustCompile(params.Get().Commands.Search)

	// Handle text commands
	switch {
	case strings.Contains(msgTxt, "kinopoisk.ru/film"):
		cmd.DownloadBest(update.Message, cmdData)

	case strings.Contains(msgTxt, "kinopoisk.ru/series"):
		cmd.ShowMovieList(update.Message, cmdData)

	case downloadRe.MatchString(msgTxt):
		cmd.ShowMetaMovieList(update.Message, cmdData, downloadRe, true)

	case searchRe.MatchString(msgTxt):
		cmd.ShowMetaMovieList(update.Message, cmdData, searchRe, false)
	}
}
