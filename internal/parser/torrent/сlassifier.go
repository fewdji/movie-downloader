package torrent

import (
	params "movie-downloader-bot/internal/config"
	"strings"
)

func (mov *Movie) SetVideoProps() {
	mov.SetResolution()
	mov.SetQuality()
	mov.SetContainer()
	mov.SetDynamicRange()
	mov.SetBitrate()
}

func (mov *Movie) SetResolution() bool {
	resolutions := params.NewParams().VideoMap.Resolution
	for _, v := range resolutions {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.Resolution = v.Name
				return true
			}
		}
	}
	return false
}

func (mov *Movie) SetQuality() bool {
	qualities := params.NewParams().VideoMap.Quality
	for _, v := range qualities {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.Quality = v.Name
				return true
			}
		}
	}
	return false
}

func (mov *Movie) SetContainer() {
	containers := params.NewParams().VideoMap.Container
	for _, v := range containers {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.Container = v.Name
				return
			}
		}
	}
	mov.Container = "AVC"
}

func (mov *Movie) SetDynamicRange() {
	ranges := params.NewParams().VideoMap.DynamicRange
	for _, v := range ranges {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.DynamicRange = v.Name
				return
			}
		}
	}
	mov.DynamicRange = "SDR"
}

func (mov *Movie) SetBitrate() bool {
	if mov.Meta.Length != 0 {
		filmLength := mov.Meta.Length / 60
		sizeMb := int(mov.Size) / 1024 * 1024
		mov.Bitrate = sizeMb / filmLength
		return true
	}
	return false
}
