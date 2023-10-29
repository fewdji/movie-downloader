package client

import (
	"fmt"
	"log"
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

func (q *Qbittorrent) Show(hash string) *Torrent {
	torrents := *q.List()
	if len(torrents) == 0 {
		log.Println("No torrents!")
		return nil
	}
	for k, t := range torrents {
		if t.Hash[0:16] == hash {
			return &torrents[k]
		}
	}
	return nil
}

func (q *Qbittorrent) List() *Torrents {
	qtorrents, err := q.client.Torrents(qbt.TorrentsOptions{})

	log.Println(qtorrents)

	newTorrents := Torrents{}

	if err != nil {
		log.Println("Can't get torrents:", err)
		return nil
	} else {
		if len(qtorrents) > 0 {
			log.Println(len(qtorrents))
			for _, t := range qtorrents {

				eta := int(t.Eta)
				if eta == 8640000 {
					eta = 0
				}

				tor := Torrent{
					Title:    t.Name,
					Hash:     t.Hash,
					State:    t.State,
					Speed:    int(t.Dlspeed),
					Progress: float64(t.Downloaded*100) / float64(t.Size),
					Category: t.Category,
					Size:     int(t.Size),
					Seeds:    int(t.NumSeeds),
					Eta:      eta,
				}

				newTorrents = append(newTorrents, tor)

				//log.Println(t.Name, t.Hash, t.State)

				//err := qb.PauseOne(t.Hash)
				//_, err := q.client.ResumeOne(t.Hash)
				//if err != nil {
				//	log.Println(err)
				//}
			}
			return &newTorrents
		} else {
			fmt.Println("No torrents found")
			return &newTorrents
		}
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

	log.Println("Downloading")
	log.Println(movie.Link)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
