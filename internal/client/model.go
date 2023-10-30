package client

import "movie-downloader-bot/internal/parser/torrent"

type Client interface {
	Download(movie *torrent.Movie, category string) error
	List() *Torrents
	Show(hash string) *Torrent
	Pause(hash string) bool
	Resume(hash string) bool
	Delete(hash string, deleteFiles bool) bool
}

type Torrent struct {
	Hash     string  `json:"hash"`
	Title    string  `json:"title"`
	Size     int     `json:"size"`
	State    string  `json:"state"`
	Speed    int     `json:"speed"`
	Progress float64 `json:"progress"`
	Category string  `json:"category"`
	Seeds    int     `json:"seeds"`
	Eta      int     `json:"eta"`
}

type Torrents []Torrent
