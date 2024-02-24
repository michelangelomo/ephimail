package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Address string
	Port    int

	Redis *redis.Client
}

func New() *Client {
	return &Client{}
}

func (c *Client) Connect() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Address, c.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.Redis = client
	return c, nil
}
