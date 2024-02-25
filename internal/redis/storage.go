package redis

import (
	"fmt"
	"time"

	"github.com/michelangelomo/ephimail/internal"
)

type Storage interface {
	StoreEmail(to, body string) error
	RetrieveEmails(to string) (string, error)
}

func (r *RedisStorage) StoreEmail(to, body string) error {
	key := fmt.Sprintf(
		"%s:%s",
		to,
		internal.GenerateHash(
			body,
			time.Now().String(),
		),
	)

	return r.Client.Set(
		r.context,
		key,
		body,
		0,
	).Err()
}

func (r *RedisStorage) RetrieveEmails(to string) (string, error) {
	return "", nil
}
