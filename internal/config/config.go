package params

import (
	"encoding/json"
	"log"
	"os"
)

type Params struct {
	Preset struct {
		Resolution string `json:"resolution"`
		Hdr        string `json:"hdr"`
		Container  string `json:"container"`
		BitrateMin int    `json:"bitrate_min"`
		BitrateMax int    `json:"bitrate_max"`
	} `json:"preset"`
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

func NewParams() *Params {
	var pm Params
	paramFile, err := os.Open("config.json")
	defer paramFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(paramFile)
	err = jsonParser.Decode(&pm)
	if err != nil {
		log.Fatal(err)
	}
	return &pm
}
