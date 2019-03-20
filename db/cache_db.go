package db

import (
	"github.com/go-redis/redis"
)

type CacheService struct {
	*redis.Client
}

func NewCacheService(host string, password string) *CacheService {
	return &CacheService{createNewRedisClient(host, "")}

}

func createNewRedisClient(host string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	// fmt.Println(pong, err)
	return client
}
