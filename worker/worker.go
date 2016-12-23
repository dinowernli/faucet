package worker

import (
	"dinowernli.me/faucet/bazel"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"golang.org/x/net/context"
)

// Worker represents an agent in the system capable of executing builds. More
// specifically, a Worker has an implementation of the faucet worker service.
type Worker struct {
	Service *workerService
}

// New creates a new worker.
func New() *Worker {
	return &Worker{Service: &workerService{
		queue: newQueue(),
	}}
}

type workerService struct {
	queue *queue
}

func (s *workerService) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return &pb_worker.StatusResponse{Healthy: true, QueueSize: int32(s.queue.size)}, nil
}

func (s *workerService) Execute(request *pb_worker.ExecutionRequest, stream pb_worker.Worker_ExecuteServer) error {
	_ = bazel.NewClient(request.Checkout.Workspace)

	s.queue.enqueue(request)

	// TODO(dino): Actually do something.
	return nil
}
