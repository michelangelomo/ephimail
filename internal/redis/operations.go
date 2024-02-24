package redis

import (
	"context"
	"log"
)

func (c *Client) AddEmail(to string, body string) {
	log.Printf("To: %s\nBody: %s\n", to, body)

	err := c.Redis.Set(context.Background(), to, body, 0).Err()
	if err != nil {
		panic(err)
	}
}
