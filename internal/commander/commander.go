package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Commander struct {
	bot     *tgbotapi.BotAPI
	meta    meta.Parser
	torrent torrent.Parser
	params  *params.Params
	cache   *cache.Cache
}

func NewCommander(bot *tgbotapi.BotAPI, meta meta.Parser, torrent torrent.Parser) *Commander {
	return &Commander{
		bot:     bot,
		meta:    meta,
		torrent: torrent,
		params:  params.NewParams(),
		cache:   cache.New(15*time.Minute, 30*time.Minute),
	}
}

func (cmd *Commander) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Printf("Recovered from panic: %v", panicValue)
		}
	}()

	if update.CallbackQuery != nil {

		clbk := strings.Split(update.CallbackQuery.Data, "|")

		cmdName := clbk[0]

		println(cmdName)
		switch cmdName {
		case "kpid":
			if len(clbk) != 2 {
				log.Println("No kpid!")
				return
			}

			kpid, err := strconv.Atoi(clbk[1])
			if err != nil {
				log.Println("Not int!")
				return
			}

			metaMovie := cmd.meta.GetByKpId(kpid)

			res := cmd.torrent.Find(metaMovie)

			var rows [][]tgbotapi.InlineKeyboardButton
			for _, mov := range res {

				cmd.cache.Set(mov.Link, mov, cache.DefaultExpiration)

				rows = append(rows,
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(
							strings.Replace(
								strings.Replace(
									fmt.Sprintf("%s %s %s %s [%.1fG] (%d)", mov.Quality, mov.Resolution, mov.Container, mov.DynamicRange, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds),
									"AVC ", "", 1),
								"SDR ", "", 1),
							fmt.Sprintf("torrent|%s", mov.Link))))
			}

			rep := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(cmd.params.StaticText.TorrentMovieSearchTitle, len(res)))
			rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

			cmd.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))

			_, err = cmd.bot.Send(rep)
			if err != nil {
				log.Println(err)
				return
			}

		case "torrent":
			log.Println("Torrent callback!")
			if len(clbk) != 2 {
				log.Fatal("No kpid!")
				return
			}

			key := clbk[1]

			mov, found := cmd.cache.Get(key)
			if !found {
				log.Fatal("Cache not found!")
				return
			}

			m := mov.(torrent.Movie)

			cmd.bot.Send(tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))

			rep := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, m.Title)
			cmd.bot.Send(rep)

			return

		default:
			log.Println("Unknown callback!")
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
			fmt.Println("MetaFilm not found!")
			return
		}
		movie := movies[0]
		fmt.Println(movie.NameRu)
		res := cmd.torrent.Find(movie)
		for _, movie := range res {
			fmt.Println(movie.Title, "||| ", movie.Quality, " ", movie.Resolution, " ", movie.Container, " ", movie.DynamicRange, " ", movie.Bitrate)
		}

		/*
			Search movie or series command
		*/
	case searchRe.MatchString(msgTxt):
		title := string(searchRe.ReplaceAll([]byte(msgTxt), []byte("")))
		metaMovies := cmd.meta.FindByTitle(title)

		if len(metaMovies) == 0 {
			fmt.Println("MetaFilm not found!")
			rep := tgbotapi.NewMessage(update.Message.Chat.ID, cmd.params.StaticText.MetaMovieNotFound)
			rep.ReplyToMessageID = update.Message.MessageID
			cmd.bot.Send(rep)
			return
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, mov := range metaMovies {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", mov.NameRu, mov.Year), fmt.Sprintf("kpid|%d", mov.Id))))
		}

		rep := tgbotapi.NewMessage(update.Message.Chat.ID, cmd.params.StaticText.MetaMovieSearchTitle)
		rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		cmd.bot.Send(rep)

	}

}
