package scheduler

import (
	"fmt"

	"dinowernli.me/faucet/bazel"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/repository"

	"github.com/Sirupsen/logrus"
)

const (
	queueCapacity       = 100
	taskChannelCapacity = 10
)

// Scheduler is in charge of taking incoming requests (on a single worker) and
// running them.
type Scheduler interface {
	// Schedule add the supplied request to the set of tasks being scheduled. If
	// this returns without error, the task can be conisdered accepted. Updates
	// on the status of the task are sent to the caller throught the returned
	// channel.
	Schedule(request *pb_worker.ExecutionRequest) (chan pb_worker.ExecutionStatus, error)

	// QueueSize returns the number of tasks currently queued.
	QueueSize() int
}

func New(logger *logrus.Logger) Scheduler {
	return &scheduler{
		logger:     logger,
		bazel:      bazel.NewClient(logger),
		queue:      make(chan *task, queueCapacity),
		repoClient: repository.NewClient(logger),
	}
}

// scheduler implements a scheduling algorithm based on a simple FIFO queue.
type scheduler struct {
	logger     *logrus.Logger
	bazel      bazel.Client
	queue      chan *task
	repoClient repository.Client
}

func (s *scheduler) Schedule(request *pb_worker.ExecutionRequest) (chan pb_worker.ExecutionStatus, error) {
	task := &task{
		request:       request,
		statusChannel: make(chan pb_worker.ExecutionStatus, taskChannelCapacity),
	}

	select {
	case s.queue <- task:
		task.statusChannel <- pb_worker.ExecutionStatus_QUEUED
		return task.statusChannel, nil
	default:
		return nil, fmt.Errorf("The scheduler queue is full")
	}
}

func (s *scheduler) QueueSize() int {
	return len(s.queue)
}

// start kicks off the background polling and processing of queued tasks.
func (s *scheduler) start() {
	go func() {
		for task := range s.queue {
			s.execute(task)
		}
	}()
	s.logger.Infof("Started scheduling consumer loop")
}

func (s *scheduler) execute(task *task) {
	task.statusChannel <- pb_worker.ExecutionStatus_RUNNING

	// Acquire a checkout of the source tree in question.
	checkoutProto := task.request.Checkout
	checkoutRoot, err := s.repoClient.Checkout(checkoutProto)
	if err != nil {
		// TODO(dino): Have the channel return a struct in order to be able to
		// send the error to the caller.
		s.logger.Errorf("Unable to get checkout for %v, error: %v", checkoutProto, err)
		task.statusChannel <- pb_worker.ExecutionStatus_FAILED
		return
	}

	// Resolve the paths in need of building and testing.
	// TODO(dino): Actually resolve.
	paths := []string{}

	// Use Bazel to build/test the paths.
	// TODO(dino): Handle build failures.
	s.bazel.Run(checkoutRoot, paths)

	task.statusChannel <- pb_worker.ExecutionStatus_SUCCEEDED
}

type task struct {
	request       *pb_worker.ExecutionRequest
	statusChannel chan pb_worker.ExecutionStatus
}
