package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/pkg/helper"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (cmd *Commander) DownloadBest(inputMessage *tgbotapi.Message, cmdData CommandData) {

	replyToMessageID := inputMessage.MessageID
	// Downloading after metaMovie callback
	if cmdData.Command != "" {
		replyToMessageID = cmdData.RootMessageId
	}

	delLastMsgIfClbk := func() {
		if cmdData.Command != "" {
			cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
		}
	}

	sendErrorMsg := func(msgTxt string) {
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = replyToMessageID
		cmd.bot.Send(errMsg)
	}

	movieId, err := strconv.Atoi(helper.GetDigitsFromStr(cmdData.Key))
	if err != nil {
		log.Println(err)
		return
	}
	metaMovie := cmd.meta.GetByKpId(movieId)
	if metaMovie == nil {
		log.Println("DownloadBest: metaMovie not found")
		delLastMsgIfClbk()
		sendErrorMsg(params.Get().StaticText.MetaMovieNotFound)
		return
	}

	res := *cmd.torrent.Find(metaMovie).BaseFilter()

	if len(res) == 0 {
		log.Println("DownloadBest: torrents not found")
		delLastMsgIfClbk()
		sendErrorMsg("Фильм не найден на трекерах!")
		return
	}

	mov := res.GetBest()

	if mov == nil {
		log.Println("DownloadBest: torrents with this params not found")
		delLastMsgIfClbk()
		cmd.ShowMovieList(inputMessage, cmdData)
		return
	}

	err = cmd.client.Download(mov, "Фильмы")

	if err != nil {
		log.Println("DownloadBest: download client error:", err)
		delLastMsgIfClbk()
		sendErrorMsg("Ошибка торрент-клиента, не удалось загрузить!")
		return
	}

	delLastMsgIfClbk()
	err = cmd.downloadMessage(mov, inputMessage.Chat.ID, "Фильмы", 0, replyToMessageID)
	if err != nil {
		log.Println("DownloadBest: send error", err)
		sendErrorMsg("Не удалось отправить сообщение о начале загрузки!")
	}
}

func (cmd *Commander) DownloadMovie(inputMessage *tgbotapi.Message, cmdData CommandData) {
	deleteMsgs := func() {
		cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID, cmdData.MovieMessageId)
	}

	sendErrorMsg := func(msgTxt string) {
		deleteMsgs()
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = cmdData.RootMessageId
		cmd.bot.Send(errMsg)
	}

	mov := torrent.Movie{}

	err := cmd.cache.Scan(cmdData.Key, &mov)
	if err != nil {
		log.Println("DownloadMovie: bad cache, metaMovie not found", err)
		sendErrorMsg("Ошибка кэша, не удалось идентифицировать фильм!")
		return
	}

	var category string

	switch cmdData.Command {
	case "df":
		category = "Фильмы"
	case "ds", "dw":
		category = "Сериалы"
	case "dt":
		category = "Телешоу"
	default:
		log.Println("Unknown category!")
		return
	}

	err = cmd.client.Download(&mov, category)

	if err != nil {
		log.Println("DownloadMovie: download client error", err)
		sendErrorMsg("Ошибка торрент-клиента, не удалось загрузить!")
	}

	deleteMsgs()

	track := 0
	if cmdData.Command == "dw" {
		track = 1
		err = cmd.tracker.Add(&mov)
		if err != nil {
			log.Println("DownloadMovie: tracking error", err)
			track = -1
		}
	}

	err = cmd.downloadMessage(&mov, inputMessage.Chat.ID, category, track, cmdData.RootMessageId)
	if err != nil {
		log.Println("DownloadMovie: send error", err)
		sendErrorMsg("Не удалось отправить сообщение о начале загрузки!")
	}
}

func (cmd *Commander) ShowMetaMovieList(inputMessage *tgbotapi.Message, cmdData CommandData, searchRe *regexp.Regexp, isDownload bool) {
	sendErrorMsg := func(msgTxt string) {
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(errMsg)
	}

	title := string(searchRe.ReplaceAll([]byte(cmdData.Key), []byte("")))
	metaMovies := cmd.meta.FindByTitle(title)

	found := len(metaMovies)
	if found == 0 {
		log.Println("ShowMetaMovieList: metaMovies not found!")
		sendErrorMsg(params.Get().StaticText.MetaMovieNotFound)
		return
	}

	cmdData = CommandData{
		RootMessageId: inputMessage.MessageID,
		Offset:        0,
	}

	limit := 6
	var rows [][]tgbotapi.InlineKeyboardButton
	i := 0
	for _, mov := range metaMovies {
		if isDownload && mov.Type == torrent.FILM_TYPE {
			cmdData.Command = "bd"
		} else {
			cmdData.Command = "l"
		}
		cmdData.Key = strconv.Itoa(mov.Id)
		serializedData, _ := json.Marshal(cmdData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%d)", mov.NameRu, mov.Year), string(serializedData))))
		i++
		if i == found || i > limit-1 {
			cmdData.Command = "cnl"
			serializedData, _ = json.Marshal(cmdData)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", string(serializedData))))
			break
		}
	}

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, params.Get().StaticText.MetaMovieSearchTitle)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyToMessageID = inputMessage.MessageID

	_, err := cmd.bot.Send(msg)
	if err != nil {
		log.Println("ShowMetaMovieList: can't send", err)
		sendErrorMsg("Ошибка, не удалось сформировать список фильмов!")
	}
}

func (cmd *Commander) ShowMovieList(inputMessage *tgbotapi.Message, cmdData CommandData) {
	parsedData := CommandData{
		RootMessageId: inputMessage.MessageID,
		Command:       "s",
	}
	if cmdData.Command != "" {
		parsedData.RootMessageId = cmdData.RootMessageId
		parsedData.MetaMessageId = inputMessage.MessageID
	}

	sendErrorMsg := func(msgTxt string) {
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = parsedData.RootMessageId
		cmd.bot.Send(errMsg)
	}

	delMetaMsg := func() {
		cmd.DeleteMessage(inputMessage.Chat.ID, parsedData.MetaMessageId)
	}

	movieId, err := strconv.Atoi(helper.GetDigitsFromStr(cmdData.Key))

	metaMovie := cmd.meta.GetByKpId(movieId)
	if metaMovie == nil {
		log.Println("ShowMovieList: metaMovie not found")
		delMetaMsg()
		sendErrorMsg("Ошибка кэша, связанный фильм не найден!")
		return
	}
	res := cmd.torrent.Find(metaMovie).BaseFilter()

	log.Println(res)

	if metaMovie.Type != torrent.FILM_TYPE {
		res.SortAsSeries()
	} else {
		res.SortBySizeDesc()
	}

	movs := *res

	found := len(movs)
	if found == 0 {
		log.Println("ShowMovieList: movies not found")
		delMetaMsg()
		sendErrorMsg("Фильм не найден на трекерах!")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	var cacheKey string
	limit := 10
	top := int(math.Min(float64(found), float64(cmdData.Offset+limit)))

	for i := cmdData.Offset; i < top; i++ {
		cacheKey = helper.Hash(movs[i].Link)
		err = cmd.cache.Set(cacheKey, movs[i], time.Hour)
		if err != nil {
			log.Println("ShowMovieList: cache error", err)
			delMetaMsg()
			sendErrorMsg("Ошибка кэша, невозможно сформировать список торрентов!")
			return
		}

		parsedData.Key = cacheKey
		serializedData, _ := json.Marshal(parsedData)

		season := " "
		if movs[i].SeasonInfo != "" {
			season += "{S" + movs[i].SeasonInfo + "} "
		}

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					strings.Replace(
						strings.Replace(
							fmt.Sprintf("%s %s%s%s %s [%.1fG] (%d)",
								movs[i].Quality, movs[i].Resolution, season, movs[i].Container, movs[i].DynamicRange, float64(movs[i].Size)/float64(1024*1024*1024), movs[i].Seeds),
							"AVC ", "", 1),
						"SDR ", "", 1),
					string(serializedData))))
	}

	if found > limit {
		parsedData.Command = "l"
		parsedData.Key = cmdData.Key

		var btns []tgbotapi.InlineKeyboardButton

		if cmdData.Offset >= limit {
			parsedData.Offset = cmdData.Offset - limit
			serializedData, _ := json.Marshal(parsedData)
			btns = append(btns, tgbotapi.NewInlineKeyboardButtonData("<<", string(serializedData)))
		}

		if found-cmdData.Offset-limit > 0 {
			parsedData.Offset = cmdData.Offset + limit
			serializedData, _ := json.Marshal(parsedData)
			btns = append(btns, tgbotapi.NewInlineKeyboardButtonData(">>", string(serializedData)))
		}
		rows = append(rows, btns)
	}

	parsedData.Command = "cnl"
	serializedData, _ := json.Marshal(parsedData)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", string(serializedData))))

	delMetaMsg()

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf(params.Get().StaticText.TorrentMovieSearchTitle, found, cmdData.Offset+1, top))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyToMessageID = parsedData.RootMessageId

	_, err = cmd.bot.Send(msg)
	if err != nil {
		log.Println("ShowMovieList: can't send", err)
		sendErrorMsg("Ошибка, не удалось сформировать список торрентов!")
	}

	return
}

func (cmd *Commander) ShowMovie(inputMessage *tgbotapi.Message, callbackId string, cmdData CommandData) {
	sendErrorMsg := func(msgTxt string) {
		cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = cmdData.RootMessageId
		cmd.bot.Send(errMsg)
	}

	mov := torrent.Movie{}

	err := cmd.cache.Scan(cmdData.Key, &mov)
	if err != nil {
		log.Println("ShowMovie: сache error", err)
		sendErrorMsg("Ошибка кэша, не удалось распознать фильм!")
		return
	}

	cmdData.Command = "placeholder"
	cmdData.MovieMessageId = inputMessage.MessageID
	serializedData, err := json.Marshal(cmdData)

	clb := string(serializedData)
	clbFilm := strings.Replace(clb, "placeholder", "df", 1)
	clbSeries := strings.Replace(clb, "placeholder", "ds", 1)
	clbShow := strings.Replace(clb, "placeholder", "dt", 1)
	clbWatch := strings.Replace(clb, "placeholder", "dw", 1)
	clbCancel := strings.Replace(clb, "placeholder", "cnl", 1)
	clbDel := strings.Replace(clb, "placeholder", "del", 1)

	pubDate, _ := time.Parse(time.RFC1123Z, mov.Published)
	mov.Published = pubDate.Format("02.01.2006 в 15:04")

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
	if mov.Meta.Type != torrent.FILM_TYPE && (mov.Meta.Completed == false || mov.Meta.Year == time.Now().Year()) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 в сериалы и отслеживать новые серии", clbWatch),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", clbCancel),
		tgbotapi.NewInlineKeyboardButtonData("Назад", clbDel),
	))

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ParseMode = "markdown"
	msg.ReplyToMessageID = inputMessage.MessageID

	ans := tgbotapi.NewCallback(callbackId, "")
	cmd.bot.Send(ans)

	_, err = cmd.bot.Send(msg)
	if err != nil {
		log.Println("ShowMovie: can't send", err)
		sendErrorMsg("Ошибка, не удалось сформировать фильм!")
	}
}

func (cmd *Commander) downloadMessage(mov *torrent.Movie, chatId int64, category string, track int, replyToMessageId int) error {
	msgText := fmt.Sprintf("Загружаю %s (%.2f Gb) с <TRACKER> в %s",
		mov.Title, float64(mov.Size)/float64(1024*1024*1024), category)
	msgText = strings.Replace(msgText, "<TRACKER>", fmt.Sprintf("[%s](%s)", mov.Tracker, mov.Link), 1)

	switch track {
	case 1:
		msgText += "\n\n*Новые серии будут загружаться по мере обновления торрента*"
	case -1:
		msgText += "\n\n*Не удалось настроить отслеживание новых серий!*"
	}

	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "markdown"

	if replyToMessageId != 0 {
		msg.ReplyToMessageID = replyToMessageId
	}

	_, err := cmd.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
