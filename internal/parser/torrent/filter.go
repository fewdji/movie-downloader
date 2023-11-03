package torrent

import (
	"fmt"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/pkg/helper"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func (movs *Movies) GetBest() *Movie {
	movs.BaseFilter().
		NoSeries().
		NoCollections().
		WithDefinedVideoParams().
		NoRemux().
		MinSeeds(params.Get().VideoFilter.Limit.Auto.SeedsMin).
		SizeLimits(params.Get().VideoFilter.Limit.Auto.SizeMin, params.Get().VideoFilter.Limit.Auto.SizeMax)

	var bests Movies
	presets := params.Get().Presets

	for _, preset := range presets {
		for _, mov := range *movs {
			if mov.Preset == preset {
				bests = append(bests, mov)
			}
		}
		if len(bests) != 0 {
			break
		}
	}

	if len(bests) == 0 {
		return nil
	}

	maxSizeKey := 0
	for k, _ := range bests {
		if bests[k].Size > bests[maxSizeKey].Size {
			maxSizeKey = k
		}
	}

	return &bests[maxSizeKey]
}

func (movs *Movies) MatchTitle() *Movies {
	for i, fl := 0, false; i < len(*movs); i++ {
		mov := (*movs)[i]
		nameRu, nameOrig := mov.Meta.NameRu, mov.Meta.NameOriginal
		year := strconv.Itoa(mov.Meta.Year)
		yearPlus := strconv.Itoa(mov.Meta.Year + 1)
		yearMinus := strconv.Itoa(mov.Meta.Year - 1)

		goodTitles := [][]string{
			{nameRu, year},
			{nameRu, nameOrig, yearPlus},
			{nameRu, nameOrig, yearMinus},
			{nameOrig, year},
		}

		if mov.Meta.Type != FILM_TYPE {
			goodTitles = append(goodTitles, []string{nameRu, nameOrig})
		}

		fl = false
		for _, gt := range goodTitles {
			if helper.ContainsAll(mov.Title, gt) {
				fl = true
			}
		}
		if !fl {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) BaseFilter() *Movies {
	return movs.MatchTitle().
		NoTrailers().
		NoBadQuality().
		NoBadFormats().
		NoDisks().
		NoStereo3D().
		NoOtherLanguages().
		NoSequels().
		WithDefinedVideoParams().
		MinSeeds(params.Get().VideoFilter.Limit.Manual.SeedsMin).
		SizeLimits(params.Get().VideoFilter.Limit.Manual.SizeMin, params.Get().VideoFilter.Limit.Manual.SizeMax)
}

func (movs *Movies) WithDefinedVideoParams() *Movies {
	for i := 0; i < len(*movs); i++ {
		if (*movs)[i].Resolution == "" || (*movs)[i].Quality == "" {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) MinSeeds(seedsNum int) *Movies {
	for i := 0; i < len(*movs); i++ {
		if (*movs)[i].Seeds < seedsNum {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) SizeLimits(min, max int64) *Movies {
	for i := 0; i < len(*movs); i++ {
		if (*movs)[i].Meta.Type != FILM_TYPE {
			continue
		}
		sizeMb := (*movs)[i].Size / (1024 * 1024)
		if (sizeMb < min && min != 0) || (sizeMb > max && max != 0) {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) NoSequels() *Movies {
	for i := 0; i < len(*movs); i++ {
		sequels := []string{
			(*movs)[i].Meta.NameRu + ":",
		}
		for k := 1; k < 13; k++ {
			sequels = append(sequels, fmt.Sprintf("%s %d", (*movs)[i].Meta.NameRu, k))
		}
		if helper.ContainsAny((*movs)[i].Title, sequels) {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) NoTrailers() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.Trailers)
}

func (movs *Movies) NoBadQuality() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.BadQuality)
}

func (movs *Movies) NoBadFormats() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.BadFormats)
}

func (movs *Movies) NoRemux() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.Remux)
}

func (movs *Movies) NoOtherLanguages() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.OtherLanguages)
}

func (movs *Movies) NoStereo3D() *Movies {
	exclude := params.Get().VideoFilter.Exclude
	for i := 0; i < len(*movs); i++ {
		if helper.ContainsAny((*movs)[i].Title, exclude.Stereo3D) {
			*movs = slices.Delete(*movs, i, i+1)
			i--
			continue
		}
		if !strings.Contains((*movs)[i].Meta.NameRu, "3D") && !strings.Contains((*movs)[i].Meta.NameOriginal, "3D") && strings.Contains((*movs)[i].Title, "3D") {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}

func (movs *Movies) NoCollections() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.Collections)
}

func (movs *Movies) NoSeries() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.Series)
}

func (movs *Movies) NoDisks() *Movies {
	return movs.remover(params.Get().VideoFilter.Exclude.Disks)
}

func (movs *Movies) SortBySizeDesc() *Movies {
	sort.Slice(*movs, func(i, j int) bool {
		return (*movs)[i].Size > (*movs)[j].Size
	})
	return movs
}

func (movs *Movies) SortAsSeries() *Movies {
	sort.Slice(*movs, func(i, j int) bool {
		ti := strings.Split((*movs)[i].SeasonInfo, "-")
		tj := strings.Split((*movs)[j].SeasonInfo, "-")
		si, err := strconv.Atoi(ti[0])
		if err != nil {
			si = 1000
		}
		sj, err := strconv.Atoi(tj[0])
		if err != nil {
			sj = 1000
		}
		if si == sj {
			return (*movs)[i].Size > (*movs)[j].Size
		}
		return si < sj
	})
	return movs
}

func (movs *Movies) remover(exclude []string) *Movies {
	for i := 0; i < len(*movs); i++ {
		if helper.ContainsAny((*movs)[i].Title, exclude) {
			*movs = slices.Delete(*movs, i, i+1)
			i--
		}
	}
	return movs
}
