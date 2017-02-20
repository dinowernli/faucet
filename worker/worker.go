package worker

import (
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
	scheduler scheduler.Scheduler
}

// New creates a new worker.
func New(logger *logrus.Logger) *Worker {
	return &Worker{
		scheduler: scheduler.New(logger),
	}
}

func (s *Worker) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return &pb_worker.StatusResponse{Healthy: true, QueueSize: int32(s.scheduler.QueueSize())}, nil
}

func (s *Worker) Execute(request *pb_worker.ExecutionRequest, stream pb_worker.Worker_ExecuteServer) error {
	out, err := s.scheduler.Schedule(request)
	if err != nil {
		return grpc.Errorf(codes.Internal, "Unable to schedule execution request: %v", err)
	}

	// Listen for updates from the scheduler and send them to the caller.
	for status := range out {
		response := &pb_worker.ExecutionResponse{
			ExecutionStatus: status,
		}
		stream.SendMsg(response)
	}

	return nil
}
