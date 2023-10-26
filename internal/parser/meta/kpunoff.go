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
	Count  int       `json:"searchFilmsCountResult"`
	Movies []KpMovie `json:"films"`
}

type KpMovie struct {
	KpId      int    `json:"filmId"`
	NameRu    string `json:"nameRu"`
	NameEn    string `json:"nameEn"`
	Year      string `json:"year"`
	Length    string `json:"filmLength"`
	Type      string `json:"type"`
	StartYear int    `json:"startYear"`
	EndYear   int    `json:"endYear"`
	Serial    bool   `json:"serial"`
	Completed bool   `json:"completed"`
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

	for _, kpMovie := range kpSearchResult.Movies {
		println(kpMovie.Year)
		kpMovie.setMovieLength()
		movieLength, _ := strconv.Atoi(kpMovie.Length)
		movieYear, _ := strconv.Atoi(kpMovie.Year)

		metaMovie := Movie{
			Id:           kpMovie.KpId,
			Type:         kpMovie.Type,
			NameRu:       kpMovie.NameRu,
			NameOriginal: kpMovie.NameEn,
			Year:         movieYear,
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

	kpMovie.setMovieLength()
	movieLength, _ := strconv.Atoi(kpMovie.Length)
	movieYear, _ := strconv.Atoi(kpMovie.Year)

	metaMovie = Movie{
		Id:           kpId,
		Type:         kpMovie.Type,
		NameRu:       kpMovie.NameRu,
		NameOriginal: kpMovie.NameEn,
		Year:         movieYear,
		Length:       movieLength,
		Completed:    kpMovie.Completed,
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

func (mov *KpMovie) setMovieLength() {
	hhmm := strings.Split(mov.Length, ":")
	if len(hhmm) != 2 {
		return
	}

	hh, err := strconv.Atoi(hhmm[0])
	mm, err := strconv.Atoi(hhmm[1])
	if err != nil {
		return
	}

	inMinutes := hh*60 + mm
	mov.Length = strconv.Itoa(inMinutes)
}
