package tracker

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"movie-downloader-bot/internal/parser/torrent"
	"os"
	"runtime"
	"time"
)

type Monitor struct {
}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (t *Monitor) Add(movie torrent.Movie) error {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Monitor) Check(movie torrent.Movie) *[]torrent.Movie {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Monitor) Monitor() {
	for {
		time.Sleep(time.Second * 10)
		fmt.Println("Task is working! Goroutine num:", runtime.NumGoroutine())
	}
}

func (t *Monitor) Migrate() {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s/%s",
		os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))

	conn, err := pgx.Connect(context.Background(), connectUrl)

	if err != nil {
		log.Println("Unable to connect to database:\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	return
}
