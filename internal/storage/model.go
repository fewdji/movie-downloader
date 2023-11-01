package storage

type Storage interface {
	Get() (*[]Tracked, error)
	Add(mov *Tracked) error
	Update(tr *Tracked) error
}

type Tracked struct {
	Meta    string `json:"meta"`
	Link    string `json:"link"`
	Tracker string `json:"tracker"`
	Title   string `json:"title"`
	Size    int64  `json:"size"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Status  int    `json:"status"`
}
