package meta

type Parser interface {
	FindByName(movieName string) []Movie
	GetById(id string) Movie
}

type Movie struct {
	Type         string
	Completed    bool
	NameRu       string
	NameOriginal string
	Year         string
	Length       int
}
