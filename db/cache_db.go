package db

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

type CacheService struct {
	*redis.Client
}

func NewCacheService(host string, password string) *CacheService {
	log.Println("NewCacheService ", host)
	return &CacheService{createNewRedisClient(host, "")}

}

func createNewRedisClient(host string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password, // no password set
		DB:           0,        // use default DB
		PoolSize:     1000,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	// fmt.Println(pong, err)
	return client
}
