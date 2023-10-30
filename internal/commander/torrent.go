package commands

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

func (cmd *Commander) ShowTorrentList(inputMessage *tgbotapi.Message, cmdData CommandData) {

	torrents := *cmd.client.List()

	log.Println(torrents)

	if len(torrents) == 0 {
		log.Println("ShowTorrentList: now active torrents!")
		if cmdData.Command == "" {
			rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Нет активных торрентов")
			rep.ReplyToMessageID = inputMessage.MessageID
			cmd.bot.Send(rep)
		}
		return
	}

	cmdData.Command = "t_sh"

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, torrent := range torrents {
		cmdData.Key = torrent.Hash[0:8]
		serializedData, _ := json.Marshal(cmdData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %s [%.1f Gb] - %.1f%%",
			torrentStateIcon(torrent.State), string([]rune(torrent.Title)[:20]), float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress), string(serializedData))))
	}

	cmdData.Command = "del"
	cancel, _ := json.Marshal(cmdData)
	refresh := strings.Replace(string(cancel), "del", "t_l", 1)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", string(cancel)),
		tgbotapi.NewInlineKeyboardButtonData("Обновить", refresh)))

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Список торрентов")
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	err := cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
	if err != nil {
		log.Println("can't delete:", err)
	}

	_, err = cmd.bot.Send(rep)
	if err != nil {
		log.Println(err)
		return
	}
}

func (cmd *Commander) ShowTorrent(inputMessage *tgbotapi.Message, cmdData CommandData) {

	if cmdData.Key == "" {
		log.Println("empty key")
		return
	}
	switch cmdData.Command {
	case "t_p":
		cmd.client.Pause(cmdData.Key)
		time.Sleep(time.Millisecond * 600)
	case "t_c":
		cmd.client.Resume(cmdData.Key)
		time.Sleep(time.Millisecond * 600)
	case "t_r", "t_rf":
		cmd.client.Delete(cmdData.Key, cmdData.Command == "t_rf")
		err := cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
		if err != nil {
			log.Println("can't delete:", err)
		}
		cmd.ShowTorrentList(inputMessage, cmdData)
		return
	}

	torrent := *cmd.client.Show(cmdData.Key)

	log.Println(torrent)

	if &torrent == nil {
		log.Println("ShowTorrent: bad has!")
		rep := tgbotapi.NewMessage(inputMessage.Chat.ID, "Торрент не существует")
		rep.ReplyToMessageID = inputMessage.MessageID
		cmd.bot.Send(rep)
		return
	}

	cmdData.Command = "placeholder"
	cmdData.MetaMessageId = inputMessage.MessageID
	serializedData, err := json.Marshal(cmdData)

	clb := string(serializedData)
	clbRun := strings.Replace(clb, "placeholder", "t_c", 1)
	clbPause := strings.Replace(clb, "placeholder", "t_p", 1)
	clbDel := strings.Replace(clb, "placeholder", "t_r", 1)
	clbDelFile := strings.Replace(clb, "placeholder", "t_rf", 1)
	clbUpdate := strings.Replace(clb, "placeholder", "t_sh", 1)
	clbBack := strings.Replace(clb, "placeholder", "t_l", 1)
	clbCancel := strings.Replace(clb, "placeholder", "del", 1)

	log.Println(torrent.Progress)

	repText := fmt.Sprintf("*%s*\nСостояние: %s\nРазмер: %.2f Gb\nЗагружено: %.2f%%\nСиды: %d",
		string([]rune(torrent.Title)[:20]), torrentStateIcon(torrent.State), float64(torrent.Size)/float64(1024*1024*1024), torrent.Progress, torrent.Seeds)

	if torrent.Eta != 0 {
		repText += fmt.Sprintf("\nСкорость: %.2f Mb/сек.\nОсталось: %d мин.", float64(torrent.Speed)/float64(1024*1024), torrent.Eta/60)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отмена", clbCancel),
		tgbotapi.NewInlineKeyboardButtonData("Назад", clbBack)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить торрент", clbDel),
			tgbotapi.NewInlineKeyboardButtonData("Удалить файл", clbDelFile)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Остановить", clbPause),
			tgbotapi.NewInlineKeyboardButtonData("Запустить", clbRun)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Обновить", clbUpdate)))

	rep := tgbotapi.NewMessage(inputMessage.Chat.ID, repText)
	rep.ParseMode = "markdown"
	rep.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	err = cmd.DeleteMessage(inputMessage.Chat.ID, inputMessage.MessageID)
	if err != nil {
		log.Println("can't delete:", err)
	}

	_, err = cmd.bot.Send(rep)
	if err != nil {
		return
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
