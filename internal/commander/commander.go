package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	params  params.Params
}

type CommandData struct {
	Offset int `json:"offset"`
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser, params params.Params) *Commander {
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		params:  params,
	}
}

func (c *Commander) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Printf("recovered from panic: %v", panicValue)
		}
	}()

	if update.CallbackQuery != nil {
		parsedData := CommandData{}
		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &parsedData)
		if err != nil {
			log.Fatal(err)
			return
		}
		strings.Split(update.CallbackQuery.Data, "_")
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Param is: %+v\n", parsedData))
		_, err = c.bot.Send(msg)
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
		c.Start(update.Message)
	}

	//downloadCommand :=

	msgTxt := strings.ToLower(strings.Trim(update.Message.Text, " /"))

	downloadRe := regexp.MustCompile(os.Getenv("DOWNLOAD_CMD"))

	switch {
	case strings.Contains(msgTxt, "kinopoisk.ru/film"):
		movieId, err := strconv.Atoi(msgTxt[strings.LastIndex(msgTxt, "/")+1:])
		if err != nil {
			log.Fatal(err)
		}
		movie := c.meta.GetByKpId(movieId)
		fmt.Println(movie.NameRu)
		res := c.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title)
		}
	case strings.Contains(msgTxt, "kinopoisk.ru/series"):
		movieId, err := strconv.Atoi(msgTxt[strings.LastIndex(msgTxt, "/")+1:])
		if err != nil {
			log.Fatal(err)
		}
		movie := c.meta.GetByKpId(movieId)
		fmt.Println(movie.NameRu)
		res := c.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title)
		}
	case downloadRe.MatchString(msgTxt):
		title := string(downloadRe.ReplaceAll([]byte(msgTxt), []byte("")))
		movies := c.meta.FindByTitle(title)
		if len(movies) == 0 {
			fmt.Println("Not found!")
			return
		}
		movie := movies[0]
		fmt.Println(movie.NameRu)
		res := c.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title, "||| ", movie.Quality, " ", movie.Resolution, " ", movie.Container, " ", movie.DynamicRange)
		}
	}

}
