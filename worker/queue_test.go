package worker

import (
	"testing"

	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Empty(t *testing.T) {
	q := newQueue()
	assert.True(t, q.empty())
}

func TestQueue_AddRemove(t *testing.T) {
	q := newQueue()
	q.enqueue(requestForPath("first/path"))
	q.enqueue(requestForPath("second/path"))
	assert.Equal(t, 2, q.size)

	request := q.dequeue()
	assert.Equal(t, "first/path", request.Targets[0].Path)
	assert.Equal(t, 1, q.size)

	q.enqueue(requestForPath("third/path"))
	assert.Equal(t, 2, q.size)

	request = q.dequeue()
	assert.Equal(t, "second/path", request.Targets[0].Path)
	assert.Equal(t, 1, q.size)

	request = q.dequeue()
	assert.Equal(t, "third/path", request.Targets[0].Path)
	assert.True(t, q.empty())
}

func requestForPath(path string) *pb_worker.ExecutionRequest {
	return &pb_worker.ExecutionRequest{
		Targets: []*pb_worker.ExecutionRequest_Target{
			&pb_worker.ExecutionRequest_Target{
				Path: path,
			},
		},
	}
}
