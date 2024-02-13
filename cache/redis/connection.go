package redis

import (
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

var (
	// redisPool Redis连接池
	redisPool *redis.Client
)

// Pool Redis连接池
func Pool() *redis.Client {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Errorf("Failed to get redis db: %s", err.Error())
		panic(err)
	}
	if redisPool == nil {
		redisPool = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		})
	}
	return redisPool
}
