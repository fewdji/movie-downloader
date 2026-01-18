package qbittorrent

import (
	"fmt"
	"log"
	"movie-downloader-bot/internal/client"
	"movie-downloader-bot/internal/parser/torrent"
	"os"

	"github.com/superturkey650/go-qbittorrent/qbt"
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
	err := q.client.Pause([]string{t.Hash})
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
	err := q.client.Resume([]string{t.Hash})
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
	err := q.client.Delete([]string{t.Hash}, deleteFiles)
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
			tor := client.Torrent{
				Title:    t.Name,
				Hash:     t.Hash,
				State:    t.State,
				Speed:    int(t.Dlspeed),
				Progress: t.Progress * 100,
				Category: t.Category,
				Size:     t.Size,
				Seeds:    int(t.NumSeeds),
				Eta:      eta,
			}
			newTorrents = append(newTorrents, tor)
		}
		return &newTorrents
	}

	return nil
}

func (q *Qbittorrent) Download(movie *torrent.Movie, category string) error {
	name := fmt.Sprintf("%s (%d)", movie.Meta.NameRu, movie.Meta.Year)
	sequential := true

	downloadOpts := qbt.DownloadOptions{
		Rename:             &name,
		Category:           &category,
		SequentialDownload: &sequential,
	}

	err := q.client.DownloadLinks([]string{movie.File}, downloadOpts)

	log.Println("Downloading: ", movie.Link)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
