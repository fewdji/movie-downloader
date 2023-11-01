package qbittorrent

import (
	"fmt"
	"log"
	"movie-downloader-bot/internal/client"
	"movie-downloader-bot/internal/parser/torrent"
	qbt "movie-downloader-bot/pkg/qbittorrent"
	"os"
)

type Qbittorrent struct {
	client *qbt.Client
}

func NewQbittorrent() *Qbittorrent {
	qb := qbt.NewClient(os.Getenv("QBT_HOST"))
	err := qb.Login(os.Getenv("QBT_USERNAME"), os.Getenv("QBT_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	return &Qbittorrent{
		client: qb,
	}
}

func (q *Qbittorrent) Show(hash string) *client.Torrent {
	res := q.List()
	if res == nil {
		return nil
	}
	torrents := *res
	if len(torrents) == 0 {
		return nil
	}
	for k, t := range torrents {
		if t.Hash[0:8] == hash {
			return &torrents[k]
		}
	}
	return nil
}

func (q *Qbittorrent) Pause(hash string) bool {
	t := q.Show(hash)
	if t == nil {
		return false
	}
	_, err := q.client.PauseOne(t.Hash)
	if err != nil {
		log.Println(err)
	}
	return true
}

func (q *Qbittorrent) Resume(hash string) bool {
	t := q.Show(hash)
	if t == nil {
		return false
	}
	_, err := q.client.ResumeOne(t.Hash)
	if err != nil {
		log.Println(err)
	}
	return true
}

func (q *Qbittorrent) Delete(hash string, deleteFiles bool) bool {
	t := q.Show(hash)
	if t == nil {
		return false
	}
	_, err := q.client.DeleteOne(t.Hash, deleteFiles)
	if err != nil {
		log.Println(err)
	}
	return true
}

func (q *Qbittorrent) List() *client.Torrents {
	newTorrents := client.Torrents{}
	torrents, err := q.client.Torrents(qbt.TorrentsOptions{})
	if err != nil {
		log.Println("Can't get torrents:", err)
		return nil
	}
	if len(torrents) > 0 {
		for _, t := range torrents {
			eta := int(t.Eta)
			if eta == 8640000 {
				eta = 0
			}
			progress := float64(0)
			if t.Size != 0 {
				progress = float64(t.Downloaded*100) + 0.001/float64(t.Size)
				if progress > 100 {
					progress = float64(100)
				}
			}
			tor := client.Torrent{
				Title:    t.Name,
				Hash:     t.Hash,
				State:    t.State,
				Speed:    int(t.Dlspeed),
				Progress: progress,
				Category: t.Category,
				Size:     t.Size,
				Seeds:    int(t.NumSeeds),
				Eta:      eta,
			}
			newTorrents = append(newTorrents, tor)
		}
		return &newTorrents
	} else {
		return nil
	}
}

func (q *Qbittorrent) Download(movie *torrent.Movie, category string) error {
	err := q.client.DownloadLink(
		movie.File,
		category,
		fmt.Sprintf("%s (%d)", movie.Meta.NameRu, movie.Meta.Year),
		true,
		false,
		"",
	)

	log.Println("Downloading: ", movie.Link)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
