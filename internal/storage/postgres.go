package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"runtime"
	"time"
)

const TRACKED_TABLE = "tracked"

type Postgres struct {
	db *sql.DB
}

type TorrentTable struct {
	Meta      string `json:"meta_string"`
	Link      string `json:"link"`
	Tracker   string `json:"tracker"`
	Title     string `json:"title"`
	Size      int64  `json:"size"`
	Published string `json:"published"`
	Updated   string `json:"updated"`
	Status    int    `json:"status"`
}

func NewPostgres() *Postgres {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("pgx", connectUrl)
	if err != nil {
		log.Println("Postgres connect error:", err)
	}
	return &Postgres{
		db: db,
	}
}

//func (t *Postgres) Add(movie torrent.Movie) error {
//	var greeting string
//	err := t.conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
//	if err != nil {
//		log.Println("Postgres error: ", err)
//		return err
//	}
//	return nil
//}

func (t *Postgres) Check(movie torrent.Movie) *[]torrent.Movie {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Postgres) Monitor() error {
	insertQuery := `INSERT
INTO tracked (meta, link, tracker, title, "size", published, updated, status)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := t.db.Exec(insertQuery,
		`{"size":1}`,
		`https://rutracker.net/forum/viewtopic.php?t=111`,
		`rutracker`,
		`Контейнер 2021`,
		494602364,
		`2022-10-26 13:22:34`,
		`2022-10-26 14:00:01`,
		0,
	)
	if err != nil {
		log.Println("Postgres insert error:", err)
		return err
	}
	return nil
}

func (t *Postgres) CreateSchema() error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS ` + TRACKED_TABLE + `
	(
	   meta json,
	   link character varying NOT NULL,
	   tracker character varying,
	   title character varying,
	   size bigint NOT NULL,
	   published timestamp without time zone,
	   updated timestamp without time zone,
	   status integer,
	   PRIMARY KEY (link)
	);`

	ownerToQuery := `ALTER TABLE IF EXISTS ` + TRACKED_TABLE + `
	   OWNER to ` + os.Getenv("POSTGRES_USERNAME") + `;`

	_, err := t.db.Exec(createTableQuery)
	if err != nil {
		log.Println("Postgres create table error:", err)
		return err
	}

	_, err = t.db.Exec(ownerToQuery)
	if err != nil {
		log.Println("Postgres owner error:", err)
		return err
	}
	return nil
}
