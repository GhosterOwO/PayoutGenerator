package util

import (
	"log"

	"github.com/go-redis/redis"
)

type redisWriter struct {
	rdb     *redis.Client
	listKey string
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.rdb.RPush(w.listKey, p).Result()
	if err != nil {
		log.Print(err)
	}
	return int(n), err
}

func NewRedisWriter(address string, key string) *redisWriter {
	return &redisWriter{
		rdb: redis.NewClient(&redis.Options{
			Addr: address,
		}),
		listKey: key,
	}
}
