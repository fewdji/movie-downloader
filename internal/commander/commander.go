package commands

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"movie-downloader-bot/internal/cache"
	"movie-downloader-bot/internal/client"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/internal/tracker"
	"movie-downloader-bot/pkg/helper"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	client  client.Client
	tracker *tracker.Tracker
	cache   cache.Cache
}

type CommandData struct {
	RootMessageId  int    `json:"z"`
	MetaMessageId  int    `json:"y"`
	MovieMessageId int    `json:"x"`
	Command        string `json:"c"`
	Key            string `json:"k"`
	Offset         int    `json:"o"`
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser, client client.Client, tracker *tracker.Tracker, cache cache.Cache) *Commander {
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		client:  client,
		tracker: tracker,
		cache:   cache,
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
		case "del":
			log.Println("del callback")
			cmd.DeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

		case "cancel":
			log.Println("cancel callback")
			cmd.DeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, cmdData.MovieMessageId, cmdData.MetaMessageId, cmdData.RootMessageId)
			cmdData = CommandData{}

		case "mm_down":
			cmd.DownloadBest(update.CallbackQuery.Message, cmdData)

		case "l":
			log.Println("Torrent callback!")
			cmd.ShowMovieList(update.CallbackQuery.Message, cmdData)

		case "s":
			log.Println("Movie show callback!")
			cmd.ShowMovie(update.CallbackQuery.Message, update.CallbackQuery.ID, cmdData)

		case "dl_f", "dl_s", "dl_t", "dl_w":
			cmd.DownloadMovie(update.CallbackQuery.Message, cmdData)

		case "t_l":
			cmd.ShowTorrentList(update.CallbackQuery.Message, cmdData)

		case "t_sh", "t_c", "t_p", "t_r", "t_rf":
			cmd.ShowTorrent(update.CallbackQuery.Message, update.CallbackQuery.ID, cmdData)

		default:
			log.Println("Unknown callback:", cmdData.Command)
		}
		return
	}

	if update.Message == nil {
		return
	}

	allowedChats := strings.Split(os.Getenv("TG_BOT_ALLOWED"), ",")
	chatId := strconv.Itoa(int(update.Message.Chat.ID))
	if !helper.ContainsAny(chatId, allowedChats) {
		log.Println("No access for:", chatId)
		return
	}

	// Handle static commands
	switch update.Message.Command() {
	case "start":
		cmd.Start(update.Message)
		return
	}

	msgTxt := strings.ToLower(update.Message.Text)
	cmdData := CommandData{Key: msgTxt}
	downloadRe := regexp.MustCompile(params.Get().Commands.Download)
	searchRe := regexp.MustCompile(params.Get().Commands.Search)
	torrentsRe := regexp.MustCompile(params.Get().Commands.Torrents)

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

		// Torrent client managment
	case torrentsRe.MatchString(msgTxt):
		cmd.ShowTorrentList(update.Message, cmdData)
	}
}
