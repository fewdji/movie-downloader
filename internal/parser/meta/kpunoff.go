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
	"time"
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
	KpId         int         `json:"filmId"`
	NameRu       string      `json:"nameRu"`
	NameEn       string      `json:"nameEn"`
	NameOriginal string      `json:"nameOriginal"`
	Year         interface{} `json:"year"`
	Length       interface{} `json:"filmLength"`
	Type         string      `json:"type"`
	StartYear    int         `json:"startYear"`
	EndYear      int         `json:"endYear"`
	Serial       bool        `json:"serial"`
	Completed    bool        `json:"completed"`
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
		log.Println(err)
		return nil
	}

	for _, kpMovie := range kpSearchResult.Movies {
		if kpMovie.Year == nil || kpMovie.NameRu == "" {
			continue
		}

		movieYear, _ := strconv.Atoi(kpMovie.Year.(string))
		metaMovie := Movie{
			Id:           kpMovie.KpId,
			Type:         kpMovie.Type,
			NameRu:       kpMovie.NameRu,
			NameOriginal: kpMovie.NameEn,
			Year:         movieYear,
		}
		metaMovies = append(metaMovies, metaMovie)
	}
	return
}

func (p *KpParser) GetByKpId(kpId int) (metaMovie *Movie) {
	apiUrl := fmt.Sprintf("%s/v2.2/films/%d", p.url, kpId)

	kpMovie := new(KpMovie)
	err := p.makeRequest(apiUrl, &kpMovie)
	if err != nil {
		log.Fatal(err)
	}

	if kpMovie.Year == nil || kpMovie.NameRu == "" {
		return nil
	}

	movieLength := 0
	switch v := kpMovie.Length.(type) {
	case float64:
		movieLength = int(v)
	}

	metaMovie = &Movie{
		Id:           kpId,
		Type:         kpMovie.Type,
		NameRu:       kpMovie.NameRu,
		NameOriginal: kpMovie.NameOriginal,
		Year:         int(kpMovie.Year.(float64)),
		Length:       movieLength,
		Completed:    kpMovie.Completed,
	}
	return
}

func (p *KpParser) makeRequest(url string, result interface{}) (err error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
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

	return nil
}
