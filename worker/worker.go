package worker

import (
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"golang.org/x/net/context"
)

type Worker struct {
	Service *workerService
}

func New() *Worker {
	return &Worker{Service: &workerService{}}
}

type workerService struct {
}

func (s *workerService) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return &pb_worker.StatusResponse{}, nil
}
