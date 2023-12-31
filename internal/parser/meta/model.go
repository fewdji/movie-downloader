package meta

type Parser interface {
	FindByTitle(movieTitle string) []Movie
	GetByKpId(kpId int) *Movie
}

type Movie struct {
	Id           int    `json:"id"`
	Type         string `json:"type"`
	Completed    bool   `json:"completed"`
	NameRu       string `json:"name_ru"`
	NameOriginal string `json:"name_original"`
	Year         int    `json:"year"`
	Length       int    `json:"length"`
}
