package torrent

import (
	"fmt"
	params "movie-downloader-bot/internal/config"
	"movie-downloader-bot/pkg/helper"
	"slices"
	"strings"
)

func (movs *Movies) BaseFilter() *Movies {
	return movs.NoTrailers().NoBadQuality().NoBadFormats().NoDisks().NoStereo3D().NoOtherLanguages().NoSequels()
}

func (movs *Movies) NoSequels() *Movies {
	for k, mov := range *movs {
		sequels := []string{
			mov.Meta.NameRu + ":",
		}
		for i := 1; i < 13; i++ {
			sequels = append(sequels, fmt.Sprintf("%s %d", mov.Meta.NameRu, i))
		}
		if helper.ContainsAny(mov.Title, sequels) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoTrailers() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Trailers) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoBadQuality() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.BadQuality) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoBadFormats() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.BadFormats) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoRemux() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Remux) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoOtherLanguages() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.OtherLanguages) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoStereo3D() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Stereo3D) {
			*movs = slices.Delete(*movs, k, k+1)
			continue
		}
		if !strings.Contains(mov.Meta.NameRu, "3D") && !strings.Contains(mov.Meta.NameOriginal, "3D") && strings.Contains(mov.Title, "3D") {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoCollections() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Collections) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoSeries() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Series) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}

func (movs *Movies) NoDisks() *Movies {
	exclude := params.NewParams().VideoFilter.Exclude
	for k, mov := range *movs {
		if helper.ContainsAny(mov.Title, exclude.Disks) {
			*movs = slices.Delete(*movs, k, k+1)
		}
	}
	return movs
}