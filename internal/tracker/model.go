package tracker

import (
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
)

type Tracker interface {
	Add(movie torrent.Movie) error
	Check() *[]Tracked
}

const TRACKED_TABLE = "tracked"

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
