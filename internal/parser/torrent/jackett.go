package torrent

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/internal/parser/meta"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type JackettParser struct {
	url    string
	token  string
	params *params.Params
}

type JackettSearchResult struct {
	XMLName       xml.Name       `xml:"rss"`
	JackettMovies []JackettMovie `xml:"channel>item"`
}

type JackettMovie struct {
	XMLName   xml.Name
	Title     string `xml:"title"`
	Tracker   string `xml:"jackettindexer"`
	Link      string `xml:"comments"`
	Published string `xml:"pubDate"`
	Size      int    `xml:"size"`
	File      string `xml:"link"`
	Props     []struct {
		XMLName xml.Name
		Name    string `xml:"name,attr"`
		Value   string `xml:"value,attr"`
	} `xml:"attr"`
	Seeds        int
	Quality      string
	Resolution   string
	DynamicRange string
	Container    string
	Bitrate      int
}

func NewJackettParser() *JackettParser {
	return &JackettParser{
		url:    os.Getenv("JACKETT_API_URL"),
		token:  os.Getenv("JACKETT_API_TOKEN"),
		params: params.NewParams(),
	}
}

func (prs *JackettParser) Find(metaMovie meta.Movie) (torrentMovies Movies) {
	q := NewQuerier(metaMovie)
	queries := q.GenerateQueries()

	var searchResult JackettSearchResult

	for _, query := range queries {
		resp, err := prs.makeRequest(query, "rutracker")
		searchResult.JackettMovies = append(searchResult.JackettMovies, resp.JackettMovies...)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, jackettMovie := range searchResult.JackettMovies {
		jackettMovie.setSeeds()

		torrentMovie := Movie{
			Meta:         metaMovie,
			Title:        jackettMovie.Title,
			Tracker:      jackettMovie.Tracker,
			Link:         jackettMovie.Link,
			Published:    jackettMovie.Published,
			Size:         jackettMovie.Size,
			File:         jackettMovie.File,
			Seeds:        jackettMovie.Seeds,
			Quality:      jackettMovie.Quality,
			Resolution:   jackettMovie.Resolution,
			DynamicRange: jackettMovie.DynamicRange,
			Container:    jackettMovie.Container,
			Bitrate:      jackettMovie.Bitrate,
		}

		torrentMovie.SetVideoProps()
		torrentMovies = append(torrentMovies, torrentMovie)
		torrentMovies.BaseFilter()
	}

	return
}

func (prs *JackettParser) makeRequest(query string, tracker string) (result *JackettSearchResult, err error) {
	apiUrl := fmt.Sprintf("%s/%s/results/torznab/api?apikey=%s&t=search&q=%s", prs.url, tracker, prs.token, url.QueryEscape(query))
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	searchResults := new(JackettSearchResult)
	err = xml.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal(err)
	}
	return searchResults, nil
}

func (mov *JackettMovie) setSeeds() {
	for _, attr := range mov.Props {
		if attr.Name == "seeders" {
			mov.Seeds, _ = strconv.Atoi(attr.Value)
		}
	}
}
