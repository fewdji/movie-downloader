package meta

type Parser interface {
	FindByTitle(movieTitle string) []Movie
	GetByKpId(kpId int) Movie
}

type Movie struct {
	Type         string
	Completed    bool
	NameRu       string
	NameOriginal string
	Year         string
	Length       int
}
