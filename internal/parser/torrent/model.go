package torrent

import (
	"encoding/json"
	"movie-downloader-bot/internal/parser/meta"
)

const (
	FILM_TYPE = "FILM"
)

type Parser interface {
	Find(meta *meta.Movie) *Movies
}

type Movies []Movie

type Movie struct {
	Meta         meta.Movie `json:"meta"`
	Title        string     `json:"title"`
	Tracker      string     `json:"tracker"`
	Link         string     `json:"link"`
	Published    string     `json:"published"`
	Size         int64      `json:"size"`
	Seeds        int        `json:"seeds"`
	Quality      string     `json:"quality"`
	Resolution   string     `json:"resolution"`
	DynamicRange string     `json:"dynamic_range"`
	Container    string     `json:"container"`
	Bitrate      int        `json:"bitrate"`
	File         string     `json:"file"`
	Preset       string     `json:"preset"`
	SeasonInfo   string     `json:"season_info"`
}

func (movie Movie) MarshalBinary() ([]byte, error) {
	return json.Marshal(movie)
}

func (movie *Movie) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &movie); err != nil {
		return err
	}
	return nil
}
