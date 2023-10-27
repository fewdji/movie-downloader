package torrent

import (
	"fmt"
	"movie-downloader-bot/internal/parser/meta"
)

type Querier struct {
	Movie meta.Movie
}

func NewQuarier(mov meta.Movie) *Querier {
	return &Querier{
		Movie: mov,
	}
}

func (q *Querier) GenerateQueries() []string {
	nameRu, nameOrig, year, movieType := q.Movie.NameRu, q.Movie.NameOriginal, q.Movie.Year, q.Movie.Type
	queries := []string{
		fmt.Sprintf("%s %d", nameRu, year),
		fmt.Sprintf("%s %s %d", nameRu, nameOrig, year+1),
		fmt.Sprintf("%s %s %d", nameRu, nameOrig, year-1),
		fmt.Sprintf("%s %d", nameOrig, year),
	}

	if movieType != "FILM" {
		queries = append(queries, fmt.Sprintf("%s %s", nameRu, nameOrig))
	}

	return queries
}
