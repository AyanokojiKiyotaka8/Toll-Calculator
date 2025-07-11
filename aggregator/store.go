package main

import (
	"fmt"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

type Storer interface {
	Insert(*types.Distance) error
	Get(int) (float64, error)
}

type MemoryStore struct {
	store map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(dist *types.Distance) error {
	s.store[dist.OBUID] += dist.Value
	return nil
}

func (s *MemoryStore) Get(id int) (float64, error) {
	val, ok := s.store[id]
	if !ok {
		return 0.0, fmt.Errorf("no recording for %d id", id)
	}
	return val, nil
}
