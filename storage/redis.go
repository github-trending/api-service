package storage

import (
	"github.com/go-redis/redis"
)

// RedisStorage represents the storage engine based on the Redis server.
type RedisStorage struct {
	Client *redis.Client
}

// NewRedisStorage returns a new redis storage client.
func NewRedisStorage(addr, auth string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	r := RedisStorage{
		Client: client,
	}

	return &r, nil
}

// HGet is A thin wrapper around redis.HGet of `github.com/go-redis/redis`.
func (r *RedisStorage) HGet(key, field string) (string, error) {
	value, err := r.Client.HGet(key, field).Result()

	return value, err
}

// HSet is A thin wrapper around redis.HSet of `github.com/go-redis/redis`.
func (r *RedisStorage) HSet(key, field, value string) (bool, error) {
	result, err := r.Client.HSet(key, field, value).Result()

	return result, err
}

// HExists is A thin wrapper around redis.HExists of `github.com/go-redis/redis`.
func (r *RedisStorage) HExists(key, field string) (bool, error) {
	result, err := r.Client.HExists(key, field).Result()

	return result, err
}
