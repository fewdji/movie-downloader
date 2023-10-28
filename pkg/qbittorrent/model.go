package qbt

// BasicTorrent holds a basic torrent object from qbittorrent
type BasicTorrent struct {
	Category               string `json:"category"`
	CompletionOn           int64  `json:"completion_on"`
	Dlspeed                int    `json:"dlspeed"`
	Eta                    int    `json:"eta"`
	ForceStart             bool   `json:"force_start"`
	Hash                   string `json:"hash"`
	Name                   string `json:"name"`
	NumComplete            int    `json:"num_complete"`
	NumIncomplete          int    `json:"num_incomplete"`
	NumLeechs              int    `json:"num_leechs"`
	NumSeeds               int    `json:"num_seeds"`
	Priority               int    `json:"priority"`
	Progress               int    `json:"progress"`
	Ratio                  int    `json:"ratio"`
	SavePath               string `json:"save_path"`
	SeqDl                  bool   `json:"seq_dl"`
	Size                   int    `json:"size"`
	State                  string `json:"state"`
	SuperSeeding           bool   `json:"super_seeding"`
	Upspeed                int    `json:"upspeed"`
	FirstLastPiecePriority bool   `json:"f_l_piece_prio"`
}

// Torrent holds a torrent object from qbittorrent
// with more information than BasicTorrent
type Torrent struct {
	AdditionDate       int     `json:"addition_date"`
	Comment            string  `json:"comment"`
	CompletionDate     int     `json:"completion_date"`
	CreatedBy          string  `json:"created_by"`
	CreationDate       int     `json:"creation_date"`
	DlLimit            int     `json:"dl_limit"`
	DlSpeed            int     `json:"dl_speed"`
	DlSpeedAvg         int     `json:"dl_speed_avg"`
	Eta                int     `json:"eta"`
	LastSeen           int     `json:"last_seen"`
	NbConnections      int     `json:"nb_connections"`
	NbConnectionsLimit int     `json:"nb_connections_limit"`
	Peers              int     `json:"peers"`
	PeersTotal         int     `json:"peers_total"`
	PieceSize          int     `json:"piece_size"`
	PiecesHave         int     `json:"pieces_have"`
	PiecesNum          int     `json:"pieces_num"`
	Reannounce         int     `json:"reannounce"`
	SavePath           string  `json:"save_path"`
	SeedingTime        int     `json:"seeding_time"`
	Seeds              int     `json:"seeds"`
	SeedsTotal         int     `json:"seeds_total"`
	ShareRatio         float64 `json:"share_ratio"`
	TimeElapsed        int     `json:"time_elapsed"`
	TotalDl            int     `json:"total_downloaded"`
	TotalDlSession     int     `json:"total_downloaded_session"`
	TotalSize          int     `json:"total_size"`
	TotalUl            int     `json:"total_uploaded"`
	TotalUlSession     int     `json:"total_uploaded_session"`
	TotalWasted        int     `json:"total_wasted"`
	UpLimit            int     `json:"up_limit"`
	UpSpeed            int     `json:"up_speed"`
	UpSpeedAvg         int     `json:"up_speed_avg"`
}

type TorrentInfo struct {
	AddedOn           int64   `json:"added_on"`
	AmountLeft        int64   `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      int64   `json:"availability"`
	Category          string  `json:"category"`
	Completed         int64   `json:"completed"`
	CompletionOn      int64   `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int64   `json:"dl_limit"`
	Dlspeed           int64   `json:"dlspeed"`
	Downloaded        int64   `json:"downloaded"`
	DownloadedSession int64   `json:"downloaded_session"`
	Eta               int64   `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int64   `json:"last_activity"`
	MagnetURI         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int64   `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int64   `json:"num_complete"`
	NumIncomplete     int64   `json:"num_incomplete"`
	NumLeechs         int64   `json:"num_leechs"`
	NumSeeds          int64   `json:"num_seeds"`
	Priority          int64   `json:"priority"`
	Progress          int64   `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        int64   `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTimeLimit  int64   `json:"seeding_time_limit"`
	SeenComplete      int64   `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int64   `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int64   `json:"time_active"`
	TotalSize         int64   `json:"total_size"`
	Tracker           string  `json:"tracker"`
	TrackersCount     int64   `json:"trackers_count"`
	UpLimit           int64   `json:"up_limit"`
	Uploaded          int64   `json:"uploaded"`
	UploadedSession   int64   `json:"uploaded_session"`
	Upspeed           int64   `json:"upspeed"`
}

// Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg           string `json:"msg"`
	NumPeers      int    `json:"num_peers"`
	NumSeeds      int    `json:"num_seeds"`
	NumLeeches    int    `json:"num_leeches"`
	NumDownloaded int    `json:"num_downloaded"`
	Tier          int    `json:"tier"`
	Status        int    `json:"status"`
	URL           string `json:"url"`
}

// WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

// TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	Index        int     `json:"index"`
	IsSeed       bool    `json:"is_seed"`
	Name         string  `json:"name"`
	Availability float32 `json:"availability"`
	Priority     int     `json:"priority"`
	Progress     int     `json:"progress"`
	Size         int     `json:"size"`
	PieceRange   []int   `json:"piece_range"`
}

type TorrentsOptions struct {
	Filter   *string  // all, downloading, completed, paused, active, inactive => optional
	Category *string  // => optional
	Sort     *string  // => optional
	Reverse  *bool    // => optional
	Limit    *int     // => optional (no negatives)
	Offset   *int     // => optional (negatives allowed)
	Hashes   []string // separated by | => optional
}

// Category of torrent
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// Categories mapping
type Categories struct {
	Category map[string]Category
}

type DownloadOptions struct {
	Savepath                   *string
	Cookie                     *string
	Category                   *string
	SkipHashChecking           *bool
	Paused                     *bool
	RootFolder                 *bool
	Rename                     *string
	UploadSpeedLimit           *int
	DownloadSpeedLimit         *int
	SequentialDownload         *bool
	AutomaticTorrentManagement *bool
	FirstLastPiecePriority     *bool
}

type PriorityValues int

const (
	Do_not_download  PriorityValues = 0
	Normal_priority  PriorityValues = 1
	High_priority    PriorityValues = 6
	Maximal_priority PriorityValues = 7
)
