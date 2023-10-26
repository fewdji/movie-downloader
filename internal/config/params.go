package params

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Params struct {
	StaticText struct {
		StartMsgTxt string `json:"start_msg"`
	} `json:"static_text"`
	Commands struct {
		Download string `json:"download"`
		Search   string `json:"search"`
	} `json:"commands"`
	Preset struct {
		Resolution string `json:"resolution"`
		Hdr        string `json:"hdr"`
		Container  string `json:"container"`
		BitrateMin int    `json:"bitrate_min"`
		BitrateMax int    `json:"bitrate_max"`
	} `json:"preset"`
	VideoFilter struct {
		Keywords struct {
			Exclude            []string `json:"exclude"`
			ExcludeCollections []string `json:"exclude_collections"`
		} `json:"keywords"`
	} `json:"video_filter"`
	VideoMap struct {
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

func NewParams() *Params {
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
