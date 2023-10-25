package torrent

import "movie-downloader-bot/internal/parser/meta"

type Parser interface {
	Find(meta.Movie) []Movie
	GetById(id string) Movie
}

type Movie struct {
	Meta         meta.Movie
	Title        string
	Tracker      string
	Link         string
	Published    string
	Size         int
	Seeds        int
	Quality      string
	Resolution   string
	DynamicRange string
	Container    string
	Bitrate      int
}