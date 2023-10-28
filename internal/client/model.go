package client

import "movie-downloader-bot/internal/parser/torrent"

type Client interface {
	Download(movie *torrent.Movie) error
}
