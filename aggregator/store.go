package main

import "github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(d types.Distance) error {
	return nil
}
