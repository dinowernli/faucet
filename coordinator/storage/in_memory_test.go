package storage

import (
	"testing"

	pb_storage "dinowernli.me/faucet/proto/storage"

	"github.com/stretchr/testify/assert"
)

func TestInMemory_NotFound(t *testing.T) {
	storage := NewInMemory()

	record, err := storage.Get("id1")
	assert.Nil(t, record)
	assert.Error(t, err)
}

func TestInMemory_PutAndGet(t *testing.T) {
	storage := NewInMemory()

	err := storage.Put(createRecord("id1"))
	assert.NoError(t, err)

	record, err := storage.Get("id1")
	assert.NoError(t, err)
	assert.Equal(t, "id1", record.Id)
}

func TestInMemory_NoAliasing(t *testing.T) {
	storage := NewInMemory()

	err := storage.Put(createRecord("id1"))
	assert.NoError(t, err)

	record, err := storage.Get("id1")
	assert.NoError(t, err)

	// We now change the content of the retrieved record and verify that the
	// storage is unaffected.
	record.Id = "some-invalid-id"
	record2, err := storage.Get("id1")
	assert.NoError(t, err)
	assert.Equal(t, "id1", record2.Id)
}

func createRecord(id string) *pb_storage.CheckRecord {
	return &pb_storage.CheckRecord{Id: id}
}
