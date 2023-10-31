package cache

import "time"

type Cache interface {
	Set(key string, val interface{}, expire time.Duration) error
	Scan(key string, val interface{}) error
	//Get(key string) (err error, val interface{})
}
