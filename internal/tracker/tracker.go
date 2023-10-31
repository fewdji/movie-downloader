package tracking

import (
	"fmt"
	"log"
	"movie-downloader-bot/internal/client"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/internal/storage"
	"runtime"
	"time"
)

//type Tracker interface {
//	Add(movie torrent.Movie) error
//	Check() *[]Tracked
//}

type Tracker struct {
	mParser meta.Parser
	tParser torrent.Parser
	client  client.Client
	storage *storage.Postgres
}

type Tracked struct {
	Meta      meta.Movie `json:"meta"`
	Link      string     `json:"link"`
	Tracker   string     `json:"tracker"`
	Title     string     `json:"title"`
	Size      int64      `json:"size"`
	Published string     `json:"published"`
	Updated   string     `json:"updated"`
	Status    int        `json:"status"`
}

func NewTracker(mParser meta.Parser, tParser torrent.Parser, client client.Client, storage *storage.Postgres) *Tracker {
	return &Tracker{
		mParser: mParser,
		tParser: tParser,
		client:  client,
		storage: storage,
	}
}

func (t *Tracker) Run() {
	for {
		err := t.storage.CreateSchema()
		//err = t.storage.Monitor()
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(time.Second * 60)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Tracker) Add(movie meta.Movie) error {

	return nil
}
