package params

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Params struct {
	StaticText struct {
		StartMsgTxt             string `json:"start_msg"`
		MetaMovieNotFound       string `json:"meta_movie_not_found"`
		MetaMovieSearchTitle    string `json:"meta_movie_search_title"`
		TorrentMovieSearchTitle string `json:"torrent_movie_search_title"`
	} `json:"static_text"`
	Commands struct {
		Download string `json:"download"`
		Search   string `json:"search"`
	} `json:"commands"`
	Presets     []string `json:"presets"`
	VideoFilter struct {
		Exclude struct {
			Trailers       []string `json:"trailers"`
			BadQuality     []string `json:"bad_quality"`
			BadFormats     []string `json:"bad_formats"`
			Disks          []string `json:"disks"`
			Remux          []string `json:"remux"`
			OtherLanguages []string `json:"other_languages"`
			Stereo3D       []string `json:"stereo_3d"`
			Collections    []string `json:"collections"`
			Series         []string `json:"series"`
		} `json:"exclude"`
		Limit struct {
			Auto struct {
				SizeMin  int `json:"size_min"`
				SizeMax  int `json:"size_max"`
				SeedsMin int `json:"seeds_min"`
			} `json:"auto"`
			Manual struct {
				SizeMin  int `json:"size_min"`
				SizeMax  int `json:"size_max"`
				SeedsMin int `json:"seeds_min"`
			} `json:"manual"`
		} `json:"limit"`
	} `json:"video_filter"`
	BitrateGoal int `json:"bitrate_goal"`
	VideoMap    struct {
		Resolution []struct {
			Name  string   `json:"name"`
			Masks []string `json:"masks"`
		} `json:"resolution"`
		Quality []struct {
			Name  string   `json:"name"`
			Masks []string `json:"masks"`
		} `json:"quality"`
		Container []struct {
			Name  string   `json:"name"`
			Masks []string `json:"masks"`
		} `json:"container"`
		DynamicRange []struct {
			Name  string   `json:"name"`
			Masks []string `json:"masks"`
		} `json:"dynamic_range"`
	} `json:"video_map"`
}

var instance *Params
var once sync.Once

func Get() *Params {
	once.Do(func() {
		instance = &Params{}
		paramFile, err := os.Open("configs/params.json")
		defer paramFile.Close()
		if err != nil {
			log.Fatal(err)
		}
		jsonParser := json.NewDecoder(paramFile)
		err = jsonParser.Decode(&instance)
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
