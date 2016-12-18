package worker

import (
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
	return &Worker{Service: &workerService{}}
}

type workerService struct {
}

func (s *workerService) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return &pb_worker.StatusResponse{}, nil
}
