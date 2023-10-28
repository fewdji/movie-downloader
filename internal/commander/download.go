package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"log"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/pkg/helper"
	"regexp"
	"strconv"
	"strings"
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
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.MetaMovieNotFound)
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
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Качаю %s (%.2f Gb) с %s", best.Title, float64(best.Size)/float64(1024*1024*1024), best.Tracker))
	rep.ReplyToMessageID = inputMessage.MessageID
	cmd.bot.Send(rep)
}

func (cmd *Commander) SearchOrDownloadByTitle(inputMessage *tgbotapi.Message, msgTxt string, searchRe *regexp.Regexp, isDownload bool) {
	title := string(searchRe.ReplaceAll([]byte(msgTxt), []byte("")))
	metaMovies := cmd.meta.FindByTitle(title)

	if len(metaMovies) == 0 {
		log.Println("SearchOrDownloadByTitle: metaMovies not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.MetaMovieNotFound)
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
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.MetaMovieSearchTitle)
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
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, cmd.params.StaticText.MetaMovieNotFound)
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
		cmd.cache.Set(cacheKey, mov, cache.DefaultExpiration)

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
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf(cmd.params.StaticText.TorrentMovieSearchTitle, len(res)))
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	cmd.bot.Send(tgbotapi.NewDeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID))
	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
	}

	return
}
