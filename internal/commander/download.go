package commands

import (
	"encoding/json"
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

func (cmd *Commander) DownloadBest(inputMessage *tgbotapi.Message, cmdData CommandData) {
	movieId, err := strconv.Atoi(helper.GetDigitsFromStr(cmdData.Key))
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
	res := *cmd.torrent.Find(metaMovie).BaseFilter()

	if len(res) == 0 {
		log.Println("DownloadMovieByLinkOrId: torrents not found!")
		return
	}

	best := res.GetBest()

	if best == nil {
		log.Println("DownloadMovieByLinkOrId: torrents with this params not found!")
		return
	}

	err = cmd.client.Download(best, "Фильмы")

	if err != nil {
		// TODO: msg about error
		return
	}

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Качаю %s (%.2f Gb) с %s в Фильмы", best.Title, float64(best.Size)/float64(1024*1024*1024), best.Tracker))

	rep.ReplyToMessageID = inputMessage.MessageID

	// Downloading by link
	if cmdData.Command != "" {
		err = cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
		if err != nil {
			log.Println("can't delete:", err)
		}
		rep.ReplyToMessageID = cmdData.RootMessageId
	}

	_, err = cmd.bot.Send(rep)
	if err != nil {
		log.Println("can't send:", err)
		return
	}
}

func (cmd *Commander) DownloadMovie(inputMessage *tgbotapi.Message, cmdData CommandData) {
	mov := torrent.Movie{}

	err := cmd.cache.Get(cmd.ctx, cmdData.Key).Scan(&mov)
	if err != nil {
		log.Println("DownloadTorrent: metaMovie not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Bad cache!")
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	var category string

	switch cmdData.Command {
	case "dl_f":
		category = "Фильмы"
	case "dl_s", "dl_w":
		category = "Сериалы"
	case "dl_t":
		category = "Телешоу"
	default:
		log.Println("Unknown category!")
		return
	}

	err = cmd.client.Download(&mov, category)

	if err != nil {
		// TODO: msg about error
		log.Println("Can't download")
		return
	}
	err = cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID, cmdData.MovieMessageId)

	repText := fmt.Sprintf("Качаю %s (%.2f Gb) с %s в %s", mov.Title, float64(mov.Size)/float64(1024*1024*1024), mov.Tracker, category)

	if cmdData.Command == "dl_w" {
		//TODO: add monitoring
		repText += "\n\n*Новые серии будут докачиваться по мере обновления торрента*"
	}

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	rep.ParseMode = "markdown"
	rep.ReplyToMessageID = cmdData.RootMessageId
	cmd.bot.Send(rep)
}

func (cmd *Commander) ShowMetaMovieList(inputMessage *tgbotapi.Message, cmdData CommandData, searchRe *regexp.Regexp, isDownload bool) {
	title := string(searchRe.ReplaceAll([]byte(cmdData.Key), []byte("")))
	metaMovies := cmd.meta.FindByTitle(title)

	found := len(metaMovies)
	if found == 0 {
		log.Println("SearchOrDownloadByTitle: metaMovies not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieNotFound)
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	parsedData := CommandData{
		RootMessageId: inputMessage.MessageID,
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	i := 0
	for _, mov := range metaMovies {
		if isDownload && mov.Type == torrent.FILM_TYPE {
			parsedData.Command = "mm_down"
		} else {
			parsedData.Command = "mm_tor"
		}
		parsedData.Key = strconv.Itoa(mov.Id)
		serializedData, _ := json.Marshal(parsedData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", mov.NameRu, mov.Year), string(serializedData))))
		i++
		if i == found || i > 5 {
			parsedData.Command = "cancel"
			serializedData, _ := json.Marshal(parsedData)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", string(serializedData))))
			break
		}
	}

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieSearchTitle)
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	rep.ReplyToMessageID = parsedData.RootMessageId

	cmd.bot.Send(rep)
}

func (cmd *Commander) ShowMovieList(inputMessage *tgbotapi.Message, cmdData CommandData) {
	movieId, err := strconv.Atoi(helper.GetDigitsFromStr(cmdData.Key))

	metaMovie := cmd.meta.GetByKpId(movieId)
	if metaMovie == nil {
		log.Println("SearchByLinkOrId: metaMovie not found!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieNotFound)
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}
	log.Println("SearchByLinkOrId: metaMovie found: ", metaMovie.NameRu)
	res := *cmd.torrent.Find(metaMovie).BaseFilter().SortBySizeAsc()

	found := len(res)
	if found == 0 {
		log.Println("SearchByLinkOrId: torrents not found!")
	}

	if found > 99 {
		res = res[found-99:]
	}

	parsedData := CommandData{}
	if cmdData.Command != "" {
		parsedData = CommandData{
			MetaMessageId: inputMessage.MessageID,
			RootMessageId: cmdData.RootMessageId,
			Command:       "m_sh",
		}
	} else {
		parsedData = CommandData{
			RootMessageId: inputMessage.MessageID,
			Command:       "m_sh",
		}
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var cacheKey string

	i := 0
	for _, mov := range res {

		cacheKey = helper.Hash(mov.Link)
		err := cmd.cache.SetEx(cmd.ctx, cacheKey, mov, time.Hour).Err()
		if err != nil {
			panic(err)
		}

		parsedData.Key = cacheKey
		serializedData, _ := json.Marshal(parsedData)

		//TODO: Add episode and season info
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					strings.Replace(
						strings.Replace(
							fmt.Sprintf("%s %s %s %s [%.1fG] (%d)", mov.Quality, mov.Resolution, mov.Container, mov.DynamicRange, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds),
							"AVC ", "", 1),
						"SDR ", "", 1),
					string(serializedData))))

		i++
		if i == found || i > 98 {
			parsedData.Command = "cancel"
			serializedData, _ := json.Marshal(parsedData)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", string(serializedData))))
			break
		}

	}

	// TODO: check list limits, filter the same, sort by size
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf(params.Get().StaticText.TorrentMovieSearchTitle, found))
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	rep.ReplyToMessageID = parsedData.RootMessageId

	// Downloading by link
	if cmdData.Command != "" {
		err = cmd.DeleteMessage(inputMessage.Chat.ID, parsedData.MetaMessageId)
		if err != nil {
			log.Println("can't delete:", err)
		}
	}

	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
	}

	return
}

func (cmd *Commander) ShowMovie(inputMessage *tgbotapi.Message, cmdData CommandData) {
	mov := torrent.Movie{}

	err := cmd.cache.Get(cmd.ctx, cmdData.Key).Scan(&mov)
	if err != nil {
		log.Fatal("Bad cache!")
	}

	cmdData.Key = helper.Hash(mov.Link)
	cmdData.Command = "dl_f"
	cmdData.MovieMessageId = inputMessage.MessageID
	serializedData, err := json.Marshal(cmdData)

	clbFilm := string(serializedData)
	clbSeries := strings.Replace(clbFilm, "dl_f", "dl_s", 1)
	clbShow := strings.Replace(clbFilm, "dl_f", "dl_t", 1)
	clbWatch := strings.Replace(clbFilm, "dl_f", "dl_w", 1)

	clbCancel := strings.Replace(clbFilm, "dl_f", "cancel", 1)
	clbDel := strings.Replace(clbFilm, "dl_f", "del", 1)

	date, _ := time.Parse(time.RFC1123Z, mov.Published)
	mov.Published = date.Format("02.01.2006 в 15:04")

	log.Println(clbCancel)
	//repText := fmt.Sprintf("*%s*\nРазмер: %.2f Gb\nСиды: %d\nТрекер: [%s](%s) \nДобавлен: %s", mov.Title, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds, mov.Link, mov.Tracker, mov.Published)

	repText := fmt.Sprintf("*%s*\nРазмер: %.2f Gb\nСиды: %d\nТрекер: <TRACKER>\nДобавлен: %s", mov.Title, float64(mov.Size)/float64(1024*1024*1024), mov.Seeds, mov.Published)
	repText = strings.Replace(repText, "<TRACKER>", fmt.Sprintf("[%s](%s)", mov.Tracker, mov.Link), 1)

	var rows [][]tgbotapi.InlineKeyboardButton

	if mov.Meta.Type == torrent.FILM_TYPE {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 Скачать в фильмы", clbFilm),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 в сериалы", clbSeries),
			tgbotapi.NewInlineKeyboardButtonData("💾 в телешоу", clbShow),
		))
	}

	if mov.Meta.Type != torrent.FILM_TYPE && mov.Meta.Completed == false {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 в сериалы и отслеживать новые серии", clbWatch),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", clbCancel),
		tgbotapi.NewInlineKeyboardButtonData("Назад", clbDel),
	))

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	rep.ParseMode = "markdown"

	//rep.ReplyToMessageID = cmdData.RootMessageId

	_, err = cmd.bot.Send(rep)
	if err != nil {
		log.Println(err)
	}
}
