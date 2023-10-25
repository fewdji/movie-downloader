package meta

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type KpParser struct {
	url   string
	token string
}

type KpSearchResult struct {
	Keyword                string `json:"keyword"`
	SearchFilmsCountResult int    `json:"searchFilmsCountResult"`
	Films                  []struct {
		FilmID     int    `json:"filmId"`
		NameRu     string `json:"nameRu,omitempty"`
		NameEn     string `json:"nameEn,omitempty"`
		Type       string `json:"type"`
		Year       string `json:"year"`
		FilmLength string `json:"filmLength,omitempty"`
	} `json:"films"`
}

type KpMovie struct {
	KinopoiskID              int     `json:"kinopoiskId"`
	NameRu                   string  `json:"nameRu"`
	NameEn                   string  `json:"nameEn"`
	NameOriginal             string  `json:"nameOriginal"`
	PosterURLPreview         string  `json:"posterUrlPreview"`
	RatingKinopoisk          float64 `json:"ratingKinopoisk"`
	RatingKinopoiskVoteCount int     `json:"ratingKinopoiskVoteCount"`
	RatingImdb               float64 `json:"ratingImdb"`
	RatingImdbVoteCount      int     `json:"ratingImdbVoteCount"`
	Year                     int     `json:"year"`
	FilmLength               int     `json:"filmLength"`
	Type                     string  `json:"type"`
	StartYear                int     `json:"startYear"`
	EndYear                  int     `json:"endYear"`
	Serial                   bool    `json:"serial"`
	ShortFilm                bool    `json:"shortFilm"`
	Completed                bool    `json:"completed"`
}

func NewKpParser() *KpParser {
	return &KpParser{
		url:   os.Getenv("KP_API_URL"),
		token: os.Getenv("KP_API_TOKEN"),
	}
}

func (p *KpParser) FindByName(movieName string) (movies []Movie) {
	apiUrl := fmt.Sprintf("%s/v2.1/films/search-by-keyword?keyword=%s", p.url, url.QueryEscape(movieName))

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", p.token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respResult, _ := io.ReadAll(resp.Body)
	kpSearchResult := new(KpSearchResult)

	err = json.Unmarshal(respResult, &kpSearchResult)
	if err != nil {
		log.Fatal(err)
	}

	for _, kpMovie := range kpSearchResult.Films {
		movieLength, _ := strconv.Atoi(kpMovie.FilmLength)
		movies = append(movies, Movie{
			Type:         kpMovie.Type,
			NameRu:       kpMovie.NameRu,
			NameOriginal: kpMovie.NameEn,
			Year:         kpMovie.Year,
			Length:       movieLength,
			Completed:    true,
		})
	}

	return
}

func (p *KpParser) GetById(id string) (metaMovie Movie) {
	movieId := id[strings.LastIndex(id, "/")+1:]
	apiUrl := fmt.Sprintf("%s/v2.2/films/%s", p.url, url.QueryEscape(movieId))

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", p.token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respResult, _ := io.ReadAll(resp.Body)
	kpMovie := new(KpMovie)

	//fmt.Println(string(respResult))

	err = json.Unmarshal(respResult, &kpMovie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(kpMovie.NameRu, " -l ", kpMovie.Year, " - ", kpMovie.Completed)

	metaMovie = Movie{
		Type:         kpMovie.Type,
		NameRu:       kpMovie.NameRu,
		NameOriginal: kpMovie.NameEn,
		Year:         strconv.Itoa(kpMovie.Year),
		Length:       kpMovie.FilmLength,
		Completed:    true,
	}

	return
}
