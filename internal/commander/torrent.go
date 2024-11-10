package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func (cmd *Commander) ShowTorrentList(inputMessage *tgbotapi.Message, cmdData CommandData) {
	sendErrorMsg := func(msgTxt string) {
		errMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgTxt)
		errMsg.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(errMsg)
	}

	delMsg := func() {
		cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
	}

	torrents := cmd.client.List()

	if torrents == nil {
		log.Println("ShowTorrentList: no active torrents")
		if cmdData.Command != "" {
			delMsg()
		}
		sendErrorMsg("Нет активных торрентов")
		return
	}

	cmdData.Command = "ts"
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, torrent := range *torrents {
		cmdData.Key = torrent.Hash[0:8]
		serializedData, _ := json.Marshal(cmdData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s [%.1f Gb] - %.1f%%",
				torrentStateIcon(torrent.State), string([]rune(torrent.Title)[:20]), float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress), string(serializedData))))
	}

	cmdData.Command = "del"
	cancel, _ := json.Marshal(cmdData)
	refresh := strings.Replace(string(cancel), "del", "tl", 1)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", string(cancel)),
		tgbotapi.NewInlineKeyboardButtonData("Обновить", refresh)))

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID,
		fmt.Sprintf("Активные торренты (%d):", len(*torrents)))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyToMessageID = inputMessage.MessageID

	delMsg()

	_, err := cmd.bot.Send(msg)
	if err != nil {
		log.Println("ShowTorrentList: can't send", err)
		sendErrorMsg("Ошибка, не удалось сформировать список торрентов!")
	}
}

func (cmd *Commander) ShowTorrent(inputMessage *tgbotapi.Message, callbackId string, cmdData CommandData) {
	delMsg := func() {
		cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
	}

	sendClbk := func(msgTxt string) {
		ans := tgbotapi.NewCallback(callbackId, msgTxt)
		cmd.bot.Send(ans)
	}

	if cmdData.Key == "" {
		log.Println("ShowTorrent: empty key")
		sendClbk("Торрент не найден!")
		return
	}

	torrent := cmd.client.Show(cmdData.Key)

	if torrent == nil {
		log.Println("ShowTorrent: not found")
		sendClbk("Торрент не найден!")
		return
	}

	switch cmdData.Command {
	case "tp":
		cmd.client.Pause(cmdData.Key)
		sendClbk("Торрент остановлен!")
	case "tc":
		cmd.client.Resume(cmdData.Key)
		sendClbk("Торрент запущен!")
	case "tr", "tf":
		cmd.client.Delete(cmdData.Key, cmdData.Command == "tf")
		sendClbk("Торрент удален!")
		delMsg()
		cmd.ShowTorrentList(inputMessage, cmdData)
		return
	}

	cmdData.Command = "placeholder"
	serializedData, err := json.Marshal(cmdData)

	clb := string(serializedData)
	clbRun := strings.Replace(clb, "placeholder", "tc", 1)
	clbPause := strings.Replace(clb, "placeholder", "tp", 1)
	clbDel := strings.Replace(clb, "placeholder", "tr", 1)
	clbDelFile := strings.Replace(clb, "placeholder", "tf", 1)
	clbUpdate := strings.Replace(clb, "placeholder", "ts", 1)
	clbBack := strings.Replace(clb, "placeholder", "tl", 1)
	clbCancel := strings.Replace(clb, "placeholder", "del", 1)

	msgText := fmt.Sprintf("*%s*\nСостояние: %s\nРазмер: %.2f Gb\nЗагружено: %.2f%%\nСиды: %d",
		string([]rune(torrent.Title)[:20]), torrentStateIcon(torrent.State), float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress, torrent.Seeds)

	if torrent.Eta != 0 {
		msgText += fmt.Sprintf("\nСкорость: %.2f Mb/сек.\nОсталось: %d мин.", float64(torrent.Speed)/float64(1024*1024), torrent.Eta/60)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Запустить", clbRun),

		tgbotapi.NewInlineKeyboardButtonData("Остановить", clbPause)),

		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить торрент", clbDel),
			tgbotapi.NewInlineKeyboardButtonData("Удалить файл", clbDelFile)),

		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", clbCancel),
			tgbotapi.NewInlineKeyboardButtonData("Назад", clbBack)),

		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Обновить", clbUpdate)))

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, msgText)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyToMessageID = inputMessage.MessageID

	delMsg()

	_, err = cmd.bot.Send(msg)
	if err != nil {
		log.Println("ShowTorrent: can't send", err)
		sendClbk("Ошибка, не удалось открыть торрент!")
	}
}

func torrentStateIcon(state string) string {
	torrentStates := map[string]string{
		"pausedDL":    "⏸",
		"pausedUP":    "⏸",
		"downloading": "⏬",
		"stalledDL":   "⏬",
		"uploading":   "⏫",
		"stalledUP":   "⏫",
		"checkingDL":  "🔄",
		"checkingUP":  "🔄",
		"queuedDL":    "⏳",
		"queuedUP":    "⏳",
		"metaDL":      "🚫",
		"error":       "🚫",
	}

	return torrentStates[state]
}
