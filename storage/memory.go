package storage

type values map[string]string

// MemoryStorage represents a "in memory" storage engine.
type MemoryStorage struct {
	storage map[string]values
}

// NewRedisStorage returns a new "in memory" storage client.
func NewMemoryStorage(addr, auth string) *MemoryStorage {
	client := MemoryStorage{
		storage: make(map[string]values),
	}

	return &client
}

// HGet implements Storage.HGet()
func (m *MemoryStorage) HGet(key, field string) (string, error) {
	f := m.storage[key]

	if len(f) == 0 {
		return "", nil
	}

	return f[field], nil
}

// HSet implements Storage.HSet()
func (m *MemoryStorage) HSet(key, field, value string) (bool, error) {
	f := m.storage[key]

	if len(f) == 0 {
		v := values{}

		v[field] = value

		m.storage[key] = v
	} else {
		m.storage[key][field] = value
	}

	return true, nil
}

// HExists implements Storage.HExists()
func (m *MemoryStorage) HExists(key, field string) (bool, error) {
	f := m.storage[key]

	if len(f) == 0 {
		return false, nil
	}

	if f[field] == "" {
		return false, nil
	}

	return true, nil
}
