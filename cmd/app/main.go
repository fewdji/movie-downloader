package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mp := meta.NewKpParser()

	movie := mp.GetById("https://www.kinopoisk.ru/film/1387021")
	//movie := mp.FindByName("Перевозчик")[0]

	fmt.Println(movie.NameRu)

	tp := torrent.NewJackettParser()

	res := tp.Find(movie)

	for _, mov := range res {
		fmt.Println(mov.Title)
	}

	//for _, v := range metafilms {
	//	fmt.Println(v.NameRu, " - ", v.NameOriginal, " - ", v.Year, " - ", v.Completed)
	//}

	//tgBotToken := os.Getenv("TG_BOT_TOKEN")
	//
	//bot, err := tgbotapi.NewBotAPI(tgBotToken)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//bot.Debug = true
	//
	//log.Printf("Authorized on %s", bot.Self.UserName)
	//
	//updateConfig := tgbotapi.UpdateConfig{
	//	Timeout: 60,
	//}
	//
	//updates := bot.GetUpdatesChan(updateConfig)
	//
	//for update := range updates {
	//	log.Printf("%+v\n", update)
	//
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You wrote: "+update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	_, err := bot.Send(msg)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
}
