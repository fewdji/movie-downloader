package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"movie-downloader-bot/internal/storage"
	"os"
)

const TRACKED_TABLE = "tracked"

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("pgx", connectUrl)
	if err != nil {
		log.Println("postgres connect error:", err)
	}
	return &Postgres{
		db: db,
	}
}

func (t *Postgres) Get() (*[]storage.Tracked, error) {
	err := t.CheckSchema()
	if err != nil {
		return nil, err
	}

	selectQuery := `SELECT meta, link, tracker, title, "size", created, status FROM ` + TRACKED_TABLE + ` WHERE status <> -1;`
	rows, err := t.db.Query(selectQuery)
	if err != nil {
		log.Println(err)
		return nil, errors.New("bad query")
	}

	var res []storage.Tracked
	for rows.Next() {
		tr := storage.Tracked{}
		err := rows.Scan(&tr.Meta, &tr.Link, &tr.Tracker, &tr.Title, &tr.Size, &tr.Created, &tr.Status)
		if err != nil {
			log.Println("can't recognize row data:", err)
			return nil, err
		}
		res = append(res, tr)
	}

	return &res, nil
}

func (t *Postgres) Add(tr *storage.Tracked) error {
	err := t.CheckSchema()
	if err != nil {
		return err
	}

	exists := 0
	checkMovieExists := `SELECT count(*) FROM ` + TRACKED_TABLE + ` WHERE link = $1;`
	if err = t.db.QueryRow(checkMovieExists, tr.Link).Scan(&exists); err != nil {
		log.Println(err)
		return errors.New("bad query")
	}
	if exists == 1 {
		log.Println("movie is already in tracked table")
		return nil
	}

	insertQuery := `INSERT INTO ` + TRACKED_TABLE + ` (meta, link, tracker, title, "size", created, status) VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err = t.db.Exec(insertQuery,
		tr.Meta,
		tr.Link,
		tr.Tracker,
		tr.Title,
		tr.Size,
		tr.Created,
		tr.Status,
	)
	if err != nil {
		log.Println("postgres insert error:", err)
		return err
	}
	return nil
}

func (t *Postgres) Update(tr *storage.Tracked) error {
	updateQuery := `UPDATE ` + TRACKED_TABLE + ` SET title = $1, "size" = $2, updated = $3, status = $4 WHERE link = $5`
	_, err := t.db.Exec(updateQuery,
		tr.Title,
		tr.Size,
		tr.Updated,
		tr.Status,
		tr.Link,
	)
	if err != nil {
		log.Println("postgres update error:", err)
		return err
	}
	return nil
}

func (t *Postgres) CheckSchema() error {
	exists := 0
	checkTableExists := `SELECT count(*) FROM pg_tables WHERE schemaname = 'public' AND tablename = '` + TRACKED_TABLE + `';`
	if err := t.db.QueryRow(checkTableExists).Scan(&exists); err != nil {
		log.Println(err)
		return errors.New("bad query")
	}
	if exists == 1 {
		return nil
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS ` + TRACKED_TABLE + `
	(
	   meta json,
	   link character varying NOT NULL,
	   tracker character varying,
	   title character varying,
	   size bigint NOT NULL,
	   created timestamp without time zone,
	   updated timestamp without time zone,
	   status integer,
	   PRIMARY KEY (link)
	);`

	ownerToQuery := `ALTER TABLE IF EXISTS ` + TRACKED_TABLE + `
	   OWNER to ` + os.Getenv("POSTGRES_USERNAME") + `;`

	_, err := t.db.Exec(createTableQuery)
	if err != nil {
		log.Println("postgres create table error:", err)
		return err
	}

	_, err = t.db.Exec(ownerToQuery)
	if err != nil {
		log.Println("postgres owner error:", err)
		return err
	}
	return nil
}
