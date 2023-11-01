package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Meta      string `field:"meta"`
	Link      string `field:"link"`
	Tracker   string `field:"tracker"`
	Title     string `field:"title"`
	Size      int64  `field:"size"`
	Published string `field:"published"`
	Updated   string `field:"updated"`
	Status    int    `field:"status"`
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

func (t *Postgres) Check(movie torrent.Movie) *[]torrent.Movie {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Postgres) Add(mov *torrent.Movie) error {
	err := t.CheckSchema()
	if err != nil {
		return err
	}

	insertQuery := `INSERT INTO ` + TRACKED_TABLE + ` (meta, link, tracker, title, "size", published, status)
VALUES($1, $2, $3, $4, $5, $6, $7)`

	meta, err := json.Marshal(mov)
	if err != nil {
		log.Println("Uninsertable, bad meta json:", err)
		return err
	}

	pubDate, _ := time.Parse(time.RFC1123Z, mov.Published)
	mov.Published = pubDate.Format("2006-01-02 15:04:05")

	_, err = t.db.Exec(insertQuery,
		string(meta),
		mov.Link,
		mov.Tracker,
		mov.Title,
		mov.Size,
		mov.Published,
		0,
	)
	if err != nil {
		log.Println("Postgres insert error:", err)
		return err
	}
	return nil
}

func (t *Postgres) CheckSchema() error {
	exists := false
	checkTableExists := `SELECT EXISTS (
    SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = '` + TRACKED_TABLE + `'
    );`
	if err := t.db.QueryRow(checkTableExists).Scan(&exists); err != nil {
		log.Println(err)
		return errors.New("bad query")
	}
	if exists {
		return nil
	}

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
