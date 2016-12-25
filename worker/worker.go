package worker

import (
	"dinowernli.me/faucet/bazel"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker/scheduler"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Worker represents an agent in the system capable of executing builds. More
// specifically, a Worker has an implementation of the faucet worker service.
type Worker struct {
	Service *workerService
}

// New creates a new worker.
func New(logger *logrus.Logger) *Worker {
	return &Worker{Service: &workerService{
		scheduler: scheduler.New(logger),
	}}
}

type workerService struct {
	scheduler scheduler.Scheduler
}

func (s *workerService) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return &pb_worker.StatusResponse{Healthy: true, QueueSize: int32(s.scheduler.QueueSize())}, nil
}

func (s *workerService) Execute(request *pb_worker.ExecutionRequest, stream pb_worker.Worker_ExecuteServer) error {
	_ = bazel.NewClient(request.Checkout.Workspace)

	_, err := s.scheduler.Schedule(request)
	if err != nil {
		return grpc.Errorf(codes.Internal, "Unable to schedule execution request: %v", err)
	}

	// TODO(dino): Listen on the channel for updates and propagate them to the claler stream.

	return nil
}
