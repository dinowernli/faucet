package scheduler

import (
	"testing"

	pb_worker "dinowernli.me/faucet/proto/service/worker"
	pb_workspace "dinowernli.me/faucet/proto/workspace"
	"dinowernli.me/faucet/worker/checkout"

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
	mockBazelClient.On("Run", mock.Anything, mock.Anything).Return(nil)

	mockCheckoutProvider := &mockCheckoutProvider{}
	mockCheckoutProvider.On("Get", mock.Anything).Return(&checkout.Checkout{}, nil)

	return &scheduler{
		logger:           logrus.New(),
		bazel:            mockBazelClient,
		queue:            make(chan *task, queueCapacityForTesting),
		checkoutProvider: mockCheckoutProvider,
	}
}

type mockBazelClient struct {
	mock.Mock
}

func (b *mockBazelClient) Run(rootPath string, targets []string) error {
	args := b.Called(rootPath, targets)
	return args.Error(0)
}

type mockCheckoutProvider struct {
	mock.Mock
}

func (p *mockCheckoutProvider) Get(proto *pb_workspace.Checkout) (*checkout.Checkout, error) {
	args := p.Called(proto)
	return args.Get(0).(*checkout.Checkout), args.Error(1)
}
