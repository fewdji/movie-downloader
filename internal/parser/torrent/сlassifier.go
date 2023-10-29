package torrent

import (
	params "movie-downloader-bot/internal/config"
	"strings"
)

func (mov *Movie) SetVideoProps() {
	mov.SetResolution()
	mov.SetDynamicRange()
	mov.SetContainer()
	mov.SetBitrate()
	mov.SetQuality()
}

func (mov *Movie) SetResolution() bool {
	resolutions := params.Get().VideoMap.Resolution
	for _, v := range resolutions {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.Resolution = v.Name
				mov.Preset += v.Name
				return true
			}
		}
	}
	return false
}

func (mov *Movie) SetQuality() bool {
	qualities := params.Get().VideoMap.Quality
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
	containers := params.Get().VideoMap.Container
	for _, v := range containers {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.Container = v.Name
				mov.Preset += " " + v.Name
				return
			}
		}
	}
	mov.Container = "AVC"
	mov.Preset += " " + mov.Container
}

func (mov *Movie) SetDynamicRange() {
	ranges := params.Get().VideoMap.DynamicRange
	for _, v := range ranges {
		for _, mask := range v.Masks {
			if strings.Contains(mov.Title, mask) {
				mov.DynamicRange = v.Name
				mov.Preset += " " + v.Name
				return
			}
		}
	}
	mov.DynamicRange = "SDR"
	mov.Preset += " " + mov.DynamicRange
}

func (mov *Movie) SetBitrate() bool {
	if mov.Meta.Length != 0 && mov.Size != 0 && mov.Meta.Type == FILM_TYPE {
		filmLength := float64(mov.Meta.Length) / float64(60)
		sizeMb := float64(mov.Size) / float64(1024*1024)
		mov.Bitrate = int(sizeMb / filmLength)

		if mov.Resolution == "" {
			return false
		}

		//Optimal for FHD SDR AVC HB
		calcBitrate := float64(params.Get().BitrateGoal)

		// TODO: move out ratio values to params
		switch mov.Resolution {
		case "UHD":
			calcBitrate *= 1.5
		case "HD":
			calcBitrate *= 0.5
		case "MD":
			calcBitrate *= 0.3
		case "SD":
			calcBitrate *= 0.15
		}

		if mov.Container == "HEVC" {
			calcBitrate *= 0.8
		}

		if mov.DynamicRange != "SDR" {
			calcBitrate *= 1.3
		}

		if mov.Bitrate == 0 {
			return false
		}

		diffBitrate := 100 - ((mov.Bitrate * 100) / int(calcBitrate))

		//println(mov.Title)
		//println(mov.Preset, " ", mov.Quality)
		//println("size: ", mov.Size/(1024*1024), " Mb")
		//println("bitrate: ", mov.Bitrate, " Mb/Sec")
		//println("calc: ", int(calcBitrate))
		//println("diff: ", diffBitrate)
		//println("-------------------")

		switch {
		case diffBitrate > 20:
			mov.Preset += " LB"
		case diffBitrate < -20:
			mov.Preset += " HB"
		default:
			mov.Preset += " MB"
		}

		return true
	}
	return false
}
