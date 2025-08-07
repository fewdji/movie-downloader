package jackett

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"movie-downloader-bot/internal/parser/meta"
	"movie-downloader-bot/internal/parser/torrent"
	"movie-downloader-bot/pkg/helper"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Size      int64  `xml:"size"`
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
		url:   os.Getenv("JACKETT_API_URL"),
		token: os.Getenv("JACKETT_API_TOKEN"),
	}
}

func (prs *JackettParser) Find(metaMovie *meta.Movie) *torrent.Movies {
	var searchResult JackettSearchResult

	var trackers []string
	trackersEnv := os.Getenv("JACKETT_TRACKERS")
	if trackersEnv == "" {
		trackers = []string{"all"}
	} else {
		trackers = strings.Split(trackersEnv, ",")
	}

	var wg sync.WaitGroup
	log.Println("Starting Jackett requests...")
	wg.Add(len(trackers) * 2)
	for _, tracker := range trackers {
		tracker := tracker
		searchF, searchS := prs.getSearchStrings(metaMovie)

		go func() {
			respF, err := prs.makeRequest(searchF, tracker)
			if err != nil {
				log.Println("Timeout:", err)
				wg.Done()
				return
			}
			searchResult.JackettMovies = append(searchResult.JackettMovies, respF.JackettMovies...)
			defer wg.Done()
		}()

		go func() {
			if searchS == "" {
				wg.Done()
				return
			}
			respS, err := prs.makeRequest(searchS, tracker)
			if err != nil {
				log.Println("Timeout:", err)
				wg.Done()
				return
			}
			searchResult.JackettMovies = append(searchResult.JackettMovies, respS.JackettMovies...)
			defer wg.Done()
		}()

	}
	wg.Wait()
	log.Println("...Jackett requests completed")

	torrentMovies := torrent.Movies{}

	unique := map[string]bool{}
	var uniqueKey string

	for _, jackettMovie := range searchResult.JackettMovies {
		uniqueKey = helper.GetDigitsFromStr(jackettMovie.Link)
		if uniqueKey != "" {
			unique[uniqueKey] = false
		}
	}

	for _, jackettMovie := range searchResult.JackettMovies {

		uniqueKey = helper.GetDigitsFromStr(jackettMovie.Link)
		if uniqueKey == "" || unique[uniqueKey] {
			continue
		}
		unique[uniqueKey] = true

		jackettMovie.setSeeds()

		log.Println(jackettMovie.Title)

		torrentMovie := torrent.Movie{
			Meta:         *metaMovie,
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

	}

	return &torrentMovies
}

func (prs *JackettParser) makeRequest(query string, tracker string) (result *JackettSearchResult, err error) {
	apiUrl := fmt.Sprintf("%s/%s/results/torznab/api?apikey=%s&t=search&q=%s", prs.url, tracker, prs.token, url.QueryEscape(query))
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	searchResults := new(JackettSearchResult)
	err = xml.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, err
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

func (prs *JackettParser) getSearchStrings(mov *meta.Movie) (s1, s2 string) {
	s1 = mov.NameRu

	nameOriginal := mov.NameOriginal
	if nameOriginal == "" {
		nameOriginal = mov.NameEn
	}

	s2 = nameOriginal

	if len([]rune(s1)) < 4 && len([]rune(nameOriginal)) < 4 {
		s1 = mov.NameRu + " " + strconv.Itoa(mov.Year)
		s2 = ""
		if len([]rune(nameOriginal)) != 0 {
			s2 = mov.NameRu + " " + nameOriginal
		}
	}

	log.Println("Finding on Jackett: ", s1, " / ", s2)

	return
}
