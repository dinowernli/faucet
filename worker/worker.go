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
	out, err := s.scheduler.Schedule(request)
	if err != nil {
		return grpc.Errorf(codes.Internal, "Unable to schedule execution request: %v", err)
	}

	// Listen for updates from the scheduler and send them to the caller.
	for status := range out {
		response := &pb_worker.ExecutionResponse{
			ExecutionStatus: createStatusProto(status),
		}
		stream.SendMsg(response)
	}

	return nil
}

func createStatusProto(status scheduler.TaskStatus) *pb_worker.ExecutionStatus {
	// TODO(dino): Finalize the status proto and implement this.
	return &pb_worker.ExecutionStatus{}
}
