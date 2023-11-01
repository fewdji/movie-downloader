package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type Redis struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	return &Redis{
		ctx:    context.Background(),
		client: client,
	}
}

func (r *Redis) Set(key string, val interface{}, expire time.Duration) error {
	err := r.client.SetEx(r.ctx, key, val, expire).Err()
	if err != nil {
		log.Println("ShowMovieList: cache error:", err)
		return err
	}
	return nil
}

func (r *Redis) Scan(key string, val interface{}) error {
	err := r.client.Get(r.ctx, key).Scan(val)
	if err != nil {
		return err
	}
	return nil
}

//func (*Redis) Get(key string) (err error, val interface{}) {}
