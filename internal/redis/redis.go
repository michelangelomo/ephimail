package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Address  string
	Port     int
	Client   *redis.Client
	EmailTTL time.Duration

	context context.Context
}

// NewStorage creates a new Redis storage with default values
func NewStorage() *RedisStorage {
	return &RedisStorage{
		context:  context.Background(),
		EmailTTL: DefaultEmailTTL,
	}
}

func (c *RedisStorage) Connect() {
	fmt.Println("connecting to redis on ", c.Address, c.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Address, c.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.Client = client
}
