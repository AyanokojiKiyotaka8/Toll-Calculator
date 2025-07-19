package main

import (
	"fmt"
	"sync"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type Storer interface {
	Insert(*types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	mu    sync.RWMutex
	store map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(dist *types.Distance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if dist == nil {
		return fmt.Errorf("distance is nil")
	}
	s.store[dist.OBUID] += dist.Value
	return nil
}

func (s *MemoryStore) Get(id int) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[id]
	if !ok {
		return 0.0, fmt.Errorf("no distance record found for OBUID %d", id)
	}
	return val, nil
}
