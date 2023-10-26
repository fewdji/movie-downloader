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
	SearchFilmsCountResult int `json:"searchFilmsCountResult"`
	Films                  []struct {
		FilmID     int    `json:"filmId"`
		NameRu     string `json:"nameRu"`
		NameEn     string `json:"nameEn"`
		Type       string `json:"type"`
		Year       string `json:"year"`
		FilmLength string `json:"filmLength"`
	} `json:"films"`
}

type KpMovie struct {
	KinopoiskID  int    `json:"kinopoiskId"`
	NameRu       string `json:"nameRu"`
	NameEn       string `json:"nameEn"`
	NameOriginal string `json:"nameOriginal"`
	Year         int    `json:"year"`
	FilmLength   string `json:"filmLength"`
	Type         string `json:"type"`
	StartYear    int    `json:"startYear"`
	EndYear      int    `json:"endYear"`
	Serial       bool   `json:"serial"`
	Completed    bool   `json:"completed"`
}

func NewKpParser() *KpParser {
	return &KpParser{
		url:   os.Getenv("KP_API_URL"),
		token: os.Getenv("KP_API_TOKEN"),
	}
}

func (p *KpParser) FindByTitle(movieTitle string) (metaMovies []Movie) {
	apiUrl := fmt.Sprintf("%s/v2.1/films/search-by-keyword?keyword=%s", p.url, url.QueryEscape(movieTitle))
	kpSearchResult := new(KpSearchResult)

	err := p.makeRequest(apiUrl, &kpSearchResult)
	if err != nil {
		log.Fatal(err)
	}

	for _, kpMovie := range kpSearchResult.Films {
		movieLength, _ := p.stringFilmLengthToInt(kpMovie.FilmLength)

		metaMovie := Movie{
			Type:         kpMovie.Type,
			NameRu:       kpMovie.NameRu,
			NameOriginal: kpMovie.NameEn,
			Year:         kpMovie.Year,
			Length:       movieLength,
			Completed:    true,
		}

		metaMovies = append(metaMovies, metaMovie)
	}

	return
}

func (p *KpParser) GetByKpId(kpId int) (metaMovie Movie) {
	apiUrl := fmt.Sprintf("%s/v2.2/films/%d", p.url, kpId)

	kpMovie := new(KpMovie)

	err := p.makeRequest(apiUrl, &kpMovie)
	if err != nil {
		log.Fatal(err)
	}

	movieLength, _ := p.stringFilmLengthToInt(kpMovie.FilmLength)

	metaMovie = Movie{
		Type:         kpMovie.Type,
		NameRu:       kpMovie.NameRu,
		NameOriginal: kpMovie.NameEn,
		Year:         strconv.Itoa(kpMovie.Year),
		Length:       movieLength,
	}

	return
}

func (p *KpParser) makeRequest(url string, result interface{}) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-KEY", p.token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respResult, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respResult, &result)
	if err != nil {
		return err
	}

	return
}

func (mov *KpMovie) SetMovieLength() {
	hhmm := strings.Split(mov.FilmLength, ":")
	if len(hhmm) != 2 {
		return
	}

	hh, err := strconv.Atoi(hhmm[0])
	mm, err := strconv.Atoi(hhmm[1])
	if err != nil {
		return
	}

	inMinutes := hh*60 + mm
	mov.FilmLength = strconv.Itoa(inMinutes)
}

func (p *KpParser) stringFilmLengthToInt(strTime string) (result int, err error) {
	hhmm := strings.Split(strTime, ":")
	if len(hhmm) != 2 {
		return
	}

	hh, err := strconv.Atoi(hhmm[0])
	mm, err := strconv.Atoi(hhmm[1])
	if err != nil {
		return
	}

	return hh*60 + mm, nil
}
