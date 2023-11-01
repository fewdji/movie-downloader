package tracking

import (
	"encoding/json"
	"fmt"
	"log"
	"movie-downloader-bot/internal/client"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/internal/storage"
	"time"
)

type Tracker struct {
	mParser meta.Parser
	tParser torrent.Parser
	client  client.Client
	storage storage.Storage
}

func NewTracker(mParser meta.Parser, tParser torrent.Parser, client client.Client, storage storage.Storage) *Tracker {
	return &Tracker{
		mParser: mParser,
		tParser: tParser,
		client:  client,
		storage: storage,
	}
}

func (t *Tracker) Run() {
	for {
		fmt.Println("tracker task started...")
		t.CheckAndUpdate()
		fmt.Println("...tracker task finished")
		time.Sleep(time.Minute * 15)
	}
}

func (t *Tracker) CheckAndUpdate() {
	trackeds, err := t.storage.Get()
	if err != nil {
		log.Println("can't get tracked movies", err)
		return
	}

	if len(*trackeds) == 0 {
		log.Println("no movies for tracking", err)
		return
	}

	for _, tracked := range *trackeds {
		metaMov := meta.Movie{}
		err = json.Unmarshal([]byte(tracked.Meta), &metaMov)
		if err != nil {
			log.Println("can't get meta movie:", err)
			continue
		}
		movs := t.tParser.Find(&metaMov).BaseFilter()
		if len(*movs) == 0 {
			log.Println(tracked.Title, "not found on trackers")
			continue
		}
		for _, mov := range *movs {
			if mov.Link != tracked.Link {
				continue
			}
			log.Println(tracked.Title, "found on trackers")
			if mov.Size == tracked.Size {
				log.Println(tracked.Title, "no updates")
				continue
			}
			log.Println(tracked.Title, "new version found")
			err = t.client.Download(&mov, "Сериалы")
			if err != nil {
				log.Println("client error:", err)
				continue
			}
			log.Println(tracked.Title, "downloading...")
			tracked.Size = mov.Size
			tracked.Updated = time.Now().Format("2006-01-02 15:04:05")
			tracked.Status += 1
			err = t.storage.Update(&tracked)
			if err != nil {
				log.Println("tracked update error:", err)
				continue
			}
			log.Println(tracked.Title, "updated")
		}
	}
}

func (t *Tracker) Add(mov *torrent.Movie) error {
	metaMov, err := json.Marshal(mov.Meta)
	if err != nil {
		log.Println("bad meta json:", err)
		return err
	}
	tracked := storage.Tracked{
		Meta:    string(metaMov),
		Link:    mov.Link,
		Tracker: mov.Tracker,
		Title:   mov.Title,
		Size:    mov.Size,
		Created: time.Now().Format("2006-01-02 15:04:05"),
		Status:  0,
	}
	err = t.storage.Add(&tracked)
	if err != nil {
		log.Println("can't add movie to the storage", err)
		return err
	}
	log.Println("added for tracking")
	return nil
}
