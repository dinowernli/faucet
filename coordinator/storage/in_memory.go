package storage

import (
	"fmt"

	pb_storage "dinowernli.me/faucet/proto/storage"
)

// NewInMemory returns a Storage instance backed by an in-memory map.
func NewInMemory() Storage {
	return &inMemoryStorage{map[string]*pb_storage.CheckRecord{}}
}

// inMemoryStorage is a very simple, in-memory implementation of the
// Storage interface.
type inMemoryStorage struct {
	data map[string]*pb_storage.CheckRecord
}

func (s *inMemoryStorage) Get(id string) (*pb_storage.CheckRecord, error) {
	result, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("Record with id [%s] does not exist", id)
	}

	// Make a defensive copy to avoid aliasing.
	snapshot := *result
	return &snapshot, nil
}

func (s *inMemoryStorage) Put(record *pb_storage.CheckRecord) error {
	id := record.Id
	if id == "" {
		return fmt.Errorf("Failed to store record with invalid id [%s]", id)
	}

	// Make a defensive copy to avoid aliasing.
	snapshot := *record
	s.data[id] = &snapshot
	return nil
}
