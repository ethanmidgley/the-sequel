package store

import (
	"sync"
)

type s struct {
	v  map[string]string
	mu sync.RWMutex
}

func (s *s) Set(key string, value string) {

	s.mu.Lock()
	s.v[key] = value
	s.mu.Unlock()

}

func (s *s) Get(key string) string {
	s.mu.RLock()
	value, ok := s.v[key]
	s.mu.RUnlock()

	if !ok {
		return ""
	}

	return value
}

func (s *s) GetValues() map[string]string {
	return s.v
}

func newStore() s {
	return s{v: make(map[string]string), mu: sync.RWMutex{}}
}

var Store s = newStore()
