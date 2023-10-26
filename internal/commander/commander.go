package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"regexp"
	"strconv"
	"strings"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	params  *params.Params
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser) *Commander {
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		params:  params.NewParams(),
	}
}

func (cmd *Commander) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Printf("Recovered from panic: %v", panicValue)
		}
	}()

	if update.CallbackQuery != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Param is: %+v\n", []byte(update.CallbackQuery.Data)))
		_, err := cmd.bot.Send(msg)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "start":
		cmd.Start(update.Message)
	}

	msgTxt := strings.ToLower(strings.Trim(update.Message.Text, " /"))

	downloadRe := regexp.MustCompile(cmd.params.Commands.Download)
	searchRe := regexp.MustCompile(cmd.params.Commands.Search)

	switch {
	case strings.Contains(msgTxt, "kinopoisk.ru/film"):
		movieId, err := strconv.Atoi(msgTxt[strings.LastIndex(msgTxt, "/")+1:])
		if err != nil {
			log.Fatal(err)
		}
		movie := cmd.meta.GetByKpId(movieId)
		fmt.Println(movie.NameRu)
		res := cmd.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title)
		}
	case strings.Contains(msgTxt, "kinopoisk.ru/series"):
		movieId, err := strconv.Atoi(msgTxt[strings.LastIndex(msgTxt, "/")+1:])
		if err != nil {
			log.Fatal(err)
		}
		movie := cmd.meta.GetByKpId(movieId)
		fmt.Println(movie.NameRu)
		res := cmd.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title)
		}
	case downloadRe.MatchString(msgTxt):
		title := string(downloadRe.ReplaceAll([]byte(msgTxt), []byte("")))
		movies := cmd.meta.FindByTitle(title)
		if len(movies) == 0 {
			fmt.Println("Not found!")
			return
		}
		movie := movies[0]
		fmt.Println(movie.NameRu)
		res := cmd.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title, "||| ", movie.Quality, " ", movie.Resolution, " ", movie.Container, " ", movie.DynamicRange, " ", movie.Bitrate)
		}

	case searchRe.MatchString(msgTxt):
		title := string(downloadRe.ReplaceAll([]byte(msgTxt), []byte("")))
		movies := cmd.meta.FindByTitle(title)
		if len(movies) == 0 {
			fmt.Println("Not found!")
			return
		}
		movie := movies[0]
		fmt.Println(movie.NameRu)
		res := cmd.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title, "||| ", movie.Quality, " ", movie.Resolution, " ", movie.Container, " ", movie.DynamicRange, " ", movie.Bitrate)
		}
	}

}
