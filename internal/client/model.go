package client

import "movie-downloader-bot/internal/parser/torrent"

type Client interface {
	Download(movie *torrent.Movie, category string) error
	//Show(torrent Torrent) error
}

type Torrent struct {
}
