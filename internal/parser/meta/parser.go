package meta

type Parser interface {
	FindByTitle(movieTitle string) []Movie
	GetByKpId(kpId int) Movie
}

type Movie struct {
	Id           int
	Type         string
	Completed    bool
	NameRu       string
	NameOriginal string
	Year         int
	Length       int
}
