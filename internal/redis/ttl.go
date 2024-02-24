// internal/redis/ttl.go
package redis

import (
	"fmt"
	"time"

	"github.com/michelangelomo/ephimail/internal"
)

// Default TTL values
const (
	DefaultEmailTTL = 24 * time.Hour
)

// StoreEmailWithTTL stores an email with the configured TTL
func (r *RedisStorage) StoreEmailWithTTL(to, body string) error {
	key := fmt.Sprintf(
		"%s:%s",
		to,
		internal.GenerateHash(
			body,
			time.Now().String(),
		),
	)

	// Set with expiration if TTL is configured
	if r.EmailTTL > 0 {
		return r.Client.Set(
			r.GetContext(),
			key,
			body,
			r.EmailTTL,
		).Err()
	}

	// Set without expiration
	return r.Client.Set(
		r.context,
		key,
		body,
		0,
	).Err()
}
