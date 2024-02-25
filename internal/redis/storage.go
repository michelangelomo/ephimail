package redis

import (
	"fmt"
	"time"

	"github.com/michelangelomo/ephimail/internal"
)

type Storage interface {
	StoreEmail(to, body string) error
	RetrieveEmails(to string) (map[string]string, error)
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

func (r *RedisStorage) RetrieveEmails(to string) (map[string]string, error) {
	var keys []string
	var result map[string]string

	iter := r.Client.Scan(r.context, 0, fmt.Sprintf("%s:*", to), 0).Iterator()
	for iter.Next(r.context) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return result, err
	}

	result = make(map[string]string, len(keys))
	for _, key := range keys {
		body, err := r.Client.Get(r.context, key).Result()
		if err != nil {
			continue
		}
		result[key] = body
	}
	return result, nil
}
