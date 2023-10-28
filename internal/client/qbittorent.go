package client

import (
	"fmt"
	"log"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/pkg/qbittorrent"
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

func (q *Qbittorrent) GetTorrents() error {
	torrents, err := q.client.Torrents(qbt.TorrentsOptions{})
	if err != nil {
		log.Println(err)
		return err
	} else {
		if len(torrents) > 0 {
			log.Println(len(torrents))
			for _, t := range torrents {
				log.Println(t.Name, t.Hash, t.State)
				//err := qb.PauseOne(t.Hash)
				_, err := q.client.ResumeOne(t.Hash)
				if err != nil {
					log.Println(err)
				}
			}

		} else {
			fmt.Println("No torrents found")
		}
	}
	return nil
}

func (q *Qbittorrent) Download(movie *torrent.Movie) error {
	err := q.client.DownloadLink(
		movie.File,
		"Фильмы",
		movie.Meta.NameRu,
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
