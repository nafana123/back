package cache

import (
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemoryStore struct {
	cache *cache.Cache
}

func NewMemoryStore(defaultExpiration time.Duration) *MemoryStore {
	return &MemoryStore{
		cache: cache.New(defaultExpiration, 5*time.Minute),
	}
}

func (m *MemoryStore) Set(email, code string) error {
	m.cache.Set(email, code, cache.DefaultExpiration)
	return nil
}

func (m *MemoryStore) Get(email string) (string, error) {
	val, found := m.cache.Get(email)
	if !found {
		return "", errors.New("код не найден или истёк")
	}

	return val.(string), nil
}

func (m *MemoryStore) Delete(email string) error {
	m.cache.Delete(email)
	return nil
}