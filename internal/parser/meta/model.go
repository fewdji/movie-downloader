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
	NameEn       string `json:"name_en"`
	NameOriginal string `json:"name_original"`
	Year         int    `json:"year"`
	EndYear      int    `json:"end_year"`
	Length       int    `json:"length"`
}
