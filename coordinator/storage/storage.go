package storage

import (
	pb_storage "dinowernli.me/faucet/proto/storage"
)

// Storage is a type used to store and retrieve data which should be considered
// persistent.
type Storage interface {
	// Get retrieves a stored record by id.
	Get(id string) (*pb_storage.CheckRecord, error)

	// Put stores a record.
	Put(record *pb_storage.CheckRecord) error
}
