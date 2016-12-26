package scheduler

import (
	"fmt"

	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker/bazel"

	"github.com/Sirupsen/logrus"
)

type TaskStatus int

const (
	StatusRunning  TaskStatus = iota
	StatusFinished TaskStatus = iota

	queueCapacity = 100
)

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
		logger: logger,
		bazel:  bazel.NewClient(logger),
		queue:  make(chan *task, queueCapacity),
	}
}

type scheduler struct {
	logger *logrus.Logger
	bazel  bazel.Client
	queue  chan *task
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

	paths := []string{}
	for _, target := range task.request.Targets {
		paths = append(paths, target.Path)
	}
	s.bazel.Run(paths)

	// TODO(dino): Handle build failures.

	task.statusChannel <- StatusFinished
}

type task struct {
	request       *pb_worker.ExecutionRequest
	statusChannel chan TaskStatus
}