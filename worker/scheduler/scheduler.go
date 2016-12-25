package scheduler

import (
	"fmt"
	"time"

	pb_worker "dinowernli.me/faucet/proto/service/worker"

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
		queue:  make(chan task, queueCapacity),
	}
}

type scheduler struct {
	logger *logrus.Logger
	queue  chan task
}

func (s *scheduler) Schedule(request *pb_worker.ExecutionRequest) (chan TaskStatus, error) {
	task := task{
		logger:        s.logger,
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
			task.execute()
		}
	}()
	s.logger.Infof("Started scheduling consumer loop")
}

type task struct {
	logger        *logrus.Logger
	request       *pb_worker.ExecutionRequest
	statusChannel chan TaskStatus
}

func (t *task) execute() {
	t.statusChannel <- StatusRunning

	// TODO(dino): Use a bazel client to do something useful here.
	t.logger.Infof("Fake executing build...")
	time.Sleep(time.Second * 3)

	t.statusChannel <- StatusFinished
}
