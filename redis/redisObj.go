package redis

import "github.com/go-redis/redis/v8"

func NewRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "106.15.192.99:6379",
		Password: "bettry", // no password set
		DB:       0,        // use default DB
	})
}
