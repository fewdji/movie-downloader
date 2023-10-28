package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/pkg/helper"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (cmd *Commander) DownloadByLinkOrId(inputMessage *tgbotapi.Message, msgTxt string, isId bool) {
	if !isId {
		msgTxt = strings.ToLower(strings.Trim(inputMessage.Text, " /"))
		msgTxt = msgTxt[strings.LastIndex(msgTxt, "/")+1:]
	}

	movieId, err := strconv.Atoi(msgTxt)
	if err != nil {
		log.Fatal(err)
	}
	metaMovie := cmd.meta.GetByKpId(movieId)
	if metaMovie == nil {
		log.Println("DownloadMovieByLinkOrId: metaMovie not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieNotFound)
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}
	log.Println("DownloadMovieByLinkOrId: metaMovie found: ", metaMovie.NameRu)
	res := cmd.torrent.Find(metaMovie)

	if len(res) == 0 {
		log.Println("DownloadMovieByLinkOrId: torrents not found!")
	}

	best := res.GetBest()

	if best == nil {
		log.Println("DownloadMovieByLinkOrId: torrents with this params not found!")
		return
	}

	// TODO: realise downloading torrents

	err = cmd.client.Download(best)

	if err != nil {
		// TODO: msg about error
		return
	}

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Качаю %s (%.2f Gb) с %s", best.Title, float64(best.Size)/float64(1024*1024*1024), best.Tracker))
	rep.ReplyToMessageID = inputMessage.MessageID
	cmd.bot.Send(rep)
}

func (cmd *Commander) SearchOrDownloadByTitle(inputMessage *tgbotapi.Message, msgTxt string, searchRe *regexp.Regexp, isDownload bool) {
	title := string(searchRe.ReplaceAll([]byte(msgTxt), []byte("")))
	metaMovies := cmd.meta.FindByTitle(title)

	if len(metaMovies) == 0 {
		log.Println("SearchOrDownloadByTitle: metaMovies not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieNotFound)
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	var cbData string
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, mov := range metaMovies {
		if isDownload && mov.Type == torrent.FILM_TYPE {
			cbData = "metamovie_download"
		} else {
			cbData = "metamovie_torrents"
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", mov.NameRu, mov.Year), fmt.Sprintf("%s|%d", cbData, mov.Id))))
	}
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieSearchTitle)
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	// TODO: delete msg with movie list

	cmd.bot.Send(rep)
}

func (cmd *Commander) SearchByLinkOrId(inputMessage *tgbotapi.Message, msgTxt string, isId bool) {
	if !isId {
		msgTxt = strings.ToLower(strings.Trim(inputMessage.Text, " /"))
		msgTxt = msgTxt[strings.LastIndex(msgTxt, "/")+1:]
	}
	movieId, err := strconv.Atoi(msgTxt)
	if err != nil {
		log.Fatal(err)
	}
	metaMovie := cmd.meta.GetByKpId(movieId)
	if metaMovie == nil {
		log.Println("SearchByLinkOrId: metaMovie not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieNotFound)
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}
	log.Println("SearchByLinkOrId: metaMovie found: ", metaMovie.NameRu)
	res := cmd.torrent.Find(metaMovie)

	if len(res) == 0 {
		log.Println("SearchByLinkOrId: torrents not found!")
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var cacheKey string
	for _, mov := range res {

		cacheKey = helper.Hash(mov.Link)
		log.Println("cache key:", cacheKey)
		err := cmd.cache.SetEx(cmd.ctx, cacheKey, mov, time.Hour).Err()
		if err != nil {
			panic(err)
		}

		//TODO: Add episode and season info
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					strings.Replace(
						strings.Replace(
							fmt.Sprintf("%s %s %s %s [%.1fG] (%d)", mov.Quality, mov.Resolution, mov.Container, mov.DynamicRange, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds),
							"AVC ", "", 1),
						"SDR ", "", 1),
					fmt.Sprintf("movie_show|%s", cacheKey))))
	}

	// TODO: check list limits, filter the same, sort by size
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf(params.Get().StaticText.TorrentMovieSearchTitle, len(res)))
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	cmd.bot.Send(tgbotapi.NewDeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID))
	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
	}

	return
}

func (cmd *Commander) ShowMovie(inputMessage *tgbotapi.Message, cacheKey string) {
	log.Println("cache key:", cacheKey)
	res := cmd.cache.Get(cmd.ctx, cacheKey)
	log.Println(res)
	if res == nil {
		log.Fatal("Bad cache!")
	}
	mov := torrent.Movie{}

	err := cmd.cache.Get(cmd.ctx, cacheKey).Scan(&mov)
	if err != nil {
		return
	}

	//json.Unmarshal(m, &mov)
	log.Println(mov)

	date, _ := time.Parse(time.RFC1123Z, mov.Published)
	mov.Published = date.Format("02.01.2006 в 15:04")

	//repText := fmt.Sprintf("*%s*\nРазмер: %.2f Gb\nСиды: %d\nТрекер: [%s](%s) \nДобавлен: %s", mov.Title, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds, mov.Link, mov.Tracker, mov.Published)

	repText := fmt.Sprintf("*%s*\nРазмер: %.2f Gb\nСиды: %d\nТрекер: <TRACKER>\nДобавлен: %s", mov.Title, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds, mov.Published)
	repText = strings.Replace(repText, "<TRACKER>", fmt.Sprintf("[%s](%s)", mov.Tracker, mov.Link), 1)

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	rep.ParseMode = "markdown"
	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
	}
}
