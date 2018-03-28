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
func NewStorage(addr, auth, debug string) Storage {
	var s Storage
	var err error

	if debug == "true" {
		s = NewMemoryStorage(addr, auth)
		log.Info("use `in memory` storage")
	} else {
		s, err = NewRedisStorage(addr, auth)

		if err != nil {
			log.WithError(err).Fatal("redis")
		}

		log.Info("use redis storage")
	}

	return s
}
