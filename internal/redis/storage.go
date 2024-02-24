package redis

import (
	"fmt"
)

type Storage interface {
	StoreEmail(to, body string) error
	RetrieveEmails(to string) (map[string]string, error)
}

// Override StoreEmail to use TTL
func (r *RedisStorage) StoreEmail(to, body string) error {
	return r.StoreEmailWithTTL(to, body)
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
