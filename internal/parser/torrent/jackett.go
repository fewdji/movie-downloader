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
	"strings"
)

type JackettParser struct {
	url   string
	token string
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
	Props     []struct {
		XMLName xml.Name
		Name    string `xml:"name,attr"`
		Value   string `xml:"value,attr"`
	} `xml:"attr"`
	Quality      string
	Resolution   string
	DynamicRange string
	Container    string
	Bitrate      int
}

func NewJackettParser() *JackettParser {
	return &JackettParser{
		url:   os.Getenv("JACKETT_API_URL"),
		token: os.Getenv("JACKETT_API_TOKEN"),
	}
}

func (p JackettParser) Find(metaMovie meta.Movie) (movies []Movie) {

	fmt.Println(metaMovie.NameRu + " " + metaMovie.Year)

	apiUrl := fmt.Sprintf("%s/all/results/torznab/api?apikey=%s&t=search&q=%s", p.url, p.token, url.QueryEscape(metaMovie.NameRu+" "+metaMovie.Year))

	resp, _ := http.Get(apiUrl)

	body, _ := io.ReadAll(resp.Body)

	//fmt.Println(string(body))

	searchResults := new(JackettSearchResult)

	err := xml.Unmarshal(body, &searchResults)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(len(searchResults.JackettMovies))

	for _, jackettMovie := range searchResults.JackettMovies {
		seeds := 0
		for _, attr := range jackettMovie.Props {
			if attr.Name == "seeds" {
				seeds, _ = strconv.Atoi(attr.Value)
			}
		}
		p.Classified(&jackettMovie, metaMovie)
		movies = append(movies, Movie{
			Meta:         metaMovie,
			Title:        jackettMovie.Title,
			Tracker:      jackettMovie.Tracker,
			Link:         jackettMovie.Link,
			Published:    jackettMovie.Published,
			Size:         jackettMovie.Size,
			Seeds:        seeds,
			Quality:      jackettMovie.Quality,
			Resolution:   jackettMovie.Resolution,
			DynamicRange: jackettMovie.DynamicRange,
			Container:    jackettMovie.Container,
			Bitrate:      jackettMovie.Bitrate,
		})
	}

	return
}

func (p JackettParser) GetById(id string) (movie Movie) {
	return
}

func (p JackettParser) Classified(movie *JackettMovie, meta meta.Movie) {
	d := params.NewParams()

RLOOP:
	for _, v := range d.VideoMap.Resolution {
		for _, mask := range v.Masks {
			if strings.Contains(movie.Title, mask) {
				movie.Resolution = v.Name
				break RLOOP
			}
		}
	}

QLOOP:
	for _, v := range d.VideoMap.Quality {
		for _, mask := range v.Masks {
			if strings.Contains(movie.Title, mask) {
				movie.Quality = v.Name
				break QLOOP
			}
		}
	}

CLOOP:
	for _, v := range d.VideoMap.Container {
		movie.Container = "AVC"
		for _, mask := range v.Masks {
			if strings.Contains(movie.Title, mask) {
				movie.Container = v.Name
				break CLOOP
			}
		}
	}

DLOOP:
	for _, v := range d.VideoMap.DynamicRange {
		movie.DynamicRange = "SDR"
		for _, mask := range v.Masks {
			if strings.Contains(movie.Title, mask) {
				movie.DynamicRange = v.Name
				break DLOOP
			}
		}
	}

	if meta.Length != 0 {
		filmLength := meta.Length / 60
		sizeMb := int(movie.Size) / 1048576
		movie.Bitrate = sizeMb / filmLength
	}

}
