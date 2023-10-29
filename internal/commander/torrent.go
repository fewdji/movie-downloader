package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func (cmd *Commander) ShowTorrentList(inputMessage *tgbotapi.Message) {

	torrents := *cmd.client.List()

	log.Println(torrents)

	if len(torrents) == 0 {
		log.Println("ShowTorrentList: now active torrents!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Нет активных торрентов")
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	torrentStates := map[string]string{
		"pausedDL":    "⏸",
		"pausedUP":    "⏸",
		"downloading": "⏬",
		"uploading":   "⏫",
		"stalledDL":   "⏬",
		"stalledUP":   "⏫",
		"checkingDL":  "🔄",
		"checkingUP":  "🔄",
		"queuedDL":    "⏳",
		"queuedUP":    "⏳",
		"metaDL":      "🚫",
		"error":       "🚫",
	}

	parsedData := CommandData{MessageId: inputMessage.MessageID}
	parsedData.Command = "t_sh"

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, torrent := range torrents {
		parsedData.Key = torrent.Hash[0:32]
		serializedData, _ := json.Marshal(parsedData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s [%.1f Gb] - %.1f%%",
			torrentStates[torrent.State], string([]rune(torrent.Title)[:20]), float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress), string(serializedData))))
	}
	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Список торрентов")
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	// TODO: delete msg

	_, err := cmd.bot.Send(rep)
	if err != nil {
		log.Println(err)
		return
	}
}

func (cmd *Commander) ShowTorrent(inputMessage *tgbotapi.Message, cmdData CommandData) {

	torrent := *cmd.client.Show(cmdData.Key)

	log.Println(torrent)

	if &torrent == nil {
		log.Println("ShowTorrent: bad has!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Торрент не существует")
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	torrentStates := map[string]string{
		"pausedDL":    "⏸",
		"pausedUP":    "⏸",
		"downloading": "⏬",
		"uploading":   "⏫",
		"stalledDL":   "⏬",
		"stalledUP":   "⏫",
		"checkingDL":  "🔄",
		"checkingUP":  "🔄",
		"queuedDL":    "⏳",
		"queuedUP":    "⏳",
		"metaDL":      "🚫",
		"error":       "🚫",
	}

	cmdData.MessageId = inputMessage.MessageID
	cmdData.Command = "t_p"
	serializedData, err := json.Marshal(cmdData)

	clbPause := string(serializedData)
	clbRun := strings.Replace(clbPause, "t_p", "t_c", 1)
	clbDel := strings.Replace(clbPause, "t_p", "t_r", 1)
	clbDelFile := strings.Replace(clbPause, "t_p", "t_rf", 1)

	println(clbRun, clbDel, clbDelFile, torrentStates)

	repText := fmt.Sprintf("*%s*\nСостояние: %s\nРазмер: %.2f Gb\nЗагружено: %.2f%%\nСиды: %d",
		string([]rune(torrent.Title)[:20]), torrentStates[torrent.State], float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress, torrent.Seeds)

	if torrent.Eta != 0 {
		repText += fmt.Sprintf("\nСкорость: %.2f Mb/сек.\nОсталось: %d мин.", float64(torrent.Speed)/float64(1024*1024), torrent.Eta/60)
	}
	//var rows [][]tgbotapi.InlineKeyboardButton

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	//rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	rep.ParseMode = "markdown"
	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
	}
}
