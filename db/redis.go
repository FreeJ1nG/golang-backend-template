package db

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/FreeJ1nG/backend-template/util"
	"github.com/redis/go-redis/v9"
)

func CreateRdb(config util.Config) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
			Password: config.RedisPassword,
			PoolSize: 190 * runtime.GOMAXPROCS(0),
			DB:       0,
		},
	)
}

func TestRdbConnection(rdb *redis.Client) {
	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("failed to connect to redis database: %s", err.Error())
	}
	fmt.Println("Connected to Redis Database")
}
