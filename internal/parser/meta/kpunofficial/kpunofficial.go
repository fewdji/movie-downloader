package kpunofficial

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"movie-downloader-bot/internal/cache"
	"movie-downloader-bot/internal/parser/meta"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type KpParser struct {
	url   string
	token string
	cache cache.Cache
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

func NewKpParser(cache cache.Cache) *KpParser {
	return &KpParser{
		url:   os.Getenv("KP_API_URL"),
		token: os.Getenv("KP_API_TOKEN"),
		cache: cache,
	}
}

func (kpMovie KpMovie) MarshalBinary() ([]byte, error) {
	return json.Marshal(kpMovie)
}

func (kpMovie *KpMovie) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &kpMovie); err != nil {
		return err
	}
	return nil
}

func (p *KpParser) FindByTitle(movieTitle string) (metaMovies []meta.Movie) {
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

		movieYear, err := strconv.Atoi(kpMovie.Year.(string))
		if err != nil || movieYear == 0 {
			continue
		}

		metaMovie := meta.Movie{
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

func (p *KpParser) GetByKpId(kpId int) (metaMovie *meta.Movie) {
	apiUrl := fmt.Sprintf("%s/v2.2/films/%d", p.url, kpId)

	kpMovie := new(KpMovie)

	cacheKey := "kp" + strconv.Itoa(kpId)
	err := p.cache.Scan(cacheKey, kpMovie)
	if err != nil {
		log.Println("GetByKpId: not found in cache", err)
		err = p.makeRequest(apiUrl, kpMovie)
		if err != nil {
			log.Println("GetByKpId: makeRequest error", err)
			return nil
		}
		err = p.cache.Set(cacheKey, kpMovie, time.Hour*12)
		if err != nil {
			log.Println("GetByKpId: set cache error", err)
		} else {
			log.Println("GetByKpId: saved to cache")
		}
	} else {
		log.Println("GetByKpId: got from cache")
	}

	if kpMovie.Year == nil || kpMovie.NameRu == "" {
		return nil
	}

	movieLength := 0
	switch v := kpMovie.Length.(type) {
	case float64:
		movieLength = int(v)
	}

	metaMovie = &meta.Movie{
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
