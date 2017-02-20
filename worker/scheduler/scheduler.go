package scheduler

import (
	"fmt"

	"dinowernli.me/faucet/bazel"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker/checkout"

	"github.com/Sirupsen/logrus"
)

type TaskStatus int

const (
	StatusRunning  TaskStatus = iota
	StatusFinished TaskStatus = iota
	StatusFailed   TaskStatus = iota

	queueCapacity = 100
)

// Scheduler is in charge of taking incoming requests (on a single worker) and
// running them.
type Scheduler interface {
	// Schedule add the supplied request to the set of tasks being scheduled. If
	// this returns without error, the task can be conisdered accepted. Updates
	// on the status of the task are sent to the caller throught the returned
	// channel.
	Schedule(request *pb_worker.ExecutionRequest) (chan TaskStatus, error)

	// QueueSize returns the number of tasks currently queued.
	QueueSize() int
}

func New(logger *logrus.Logger) Scheduler {
	return &scheduler{
		logger:           logger,
		bazel:            bazel.NewClient(logger),
		queue:            make(chan *task, queueCapacity),
		checkoutProvider: checkout.NewProvider(),
	}
}

// scheduler implements a scheduling algorithm based on a simple FIFO queue.
type scheduler struct {
	logger           *logrus.Logger
	bazel            bazel.Client
	queue            chan *task
	checkoutProvider checkout.CheckoutProvider
}

func (s *scheduler) Schedule(request *pb_worker.ExecutionRequest) (chan TaskStatus, error) {
	task := &task{
		request:       request,
		statusChannel: make(chan TaskStatus),
	}

	select {
	case s.queue <- task:
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
	task.statusChannel <- StatusRunning

	// Acquire a checkout of the source tree in question.
	checkoutProto := task.request.Checkout
	checkout, err := s.checkoutProvider.Get(checkoutProto)
	if err != nil {
		// TODO(dino): Have the channel return a struct in order to be able to
		// send the error to the caller.
		s.logger.Errorf("Unable to get checkout for %v, error: %v", checkoutProto, err)
		task.statusChannel <- StatusFailed
		return
	}
	defer checkout.Close()

	// Resolve the paths in need of building and testing.
	// TODO(dino): Actually resolve.
	paths := []string{}

	// Use Bazel to build/test the paths.
	// TODO(dino): Handle build failures.
	s.bazel.Run(checkout.RootPath, paths)

	task.statusChannel <- StatusFinished
}

type task struct {
	request       *pb_worker.ExecutionRequest
	statusChannel chan TaskStatus
}
