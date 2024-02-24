package redis

import (
	"context"
)

// GetContext returns the storage context
func (r *RedisStorage) GetContext() context.Context {
	return r.context
}
