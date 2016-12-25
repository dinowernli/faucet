package scheduler

import (
	"testing"

	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	queueCapacityForTesting = 2
)

func TestQueueFillsUp(t *testing.T) {
	scheduler := createScheduler()

	_, err := scheduler.Schedule(&pb_worker.ExecutionRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 1, scheduler.QueueSize())

	_, err = scheduler.Schedule(&pb_worker.ExecutionRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 2, scheduler.QueueSize())

	// The queue is now full.
	_, err = scheduler.Schedule(&pb_worker.ExecutionRequest{})
	assert.Error(t, err)
	assert.Equal(t, 2, scheduler.QueueSize())
}

func TestStatusUpdates(t *testing.T) {
	scheduler := createScheduler()
	scheduler.start()

	out, err := scheduler.Schedule(&pb_worker.ExecutionRequest{})
	assert.NoError(t, err)

	firstUpdate := <-out
	assert.Equal(t, StatusRunning, firstUpdate)
	secondUpdate := <-out
	assert.Equal(t, StatusFinished, secondUpdate)
}

func createScheduler() *scheduler {
	mockBazelClient := &mockBazelClient{}
	mockBazelClient.On("Run", mock.Anything).Return(nil)

	return &scheduler{
		logger: logrus.New(),
		bazel:  mockBazelClient,
		queue:  make(chan *task, queueCapacityForTesting),
	}
}

type mockBazelClient struct {
	mock.Mock
}

func (b *mockBazelClient) Run(targets []string) error {
	args := b.Called(targets)
	return args.Error(0)
}
