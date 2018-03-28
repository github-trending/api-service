package storage

import (
	"github.com/apex/log"
)

// Storage represents a new storage type.
type Storage interface {
	HGet(key, field string) (string, error)
	HSet(key, field, value string) (bool, error)
	HExists(key, field string) (bool, error)
}

// NewStorage returns a new storage.
func NewStorage(addr, auth string) Storage {
	s, err := NewRedisStorage(addr, auth)

	if err != nil {
		log.WithError(err).Fatal("redis")
	}

	return s
}
