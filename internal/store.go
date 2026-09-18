package main

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string // code -> original URL
	hits map[string]int    // code -> hit count
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
		hits: make(map[string]int),
	}
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[code]
	return url, ok
}

func (s *Store) Set(code, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[code] = url
}

func (s *Store) IncrementHits(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hits[code]++
}

func (s *Store) Hits(code string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hits[code]
}
