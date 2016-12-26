package coordinator

import (
	"fmt"
	"time"

	"dinowernli.me/faucet/config"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	workerPollFrequency = time.Second * 2
)

// Coordinator represents an agent in the system which implements the faucet
// coordinator service. The coordinator uses faucet workers it knows of in
// order to make sure builds get executed.
type Coordinator struct {
	Service *coordinatorService
	config  config.Config
	logger  *logrus.Logger
}

// New creates a new coordinator, and is otherwise side-effect free.
func New(logger *logrus.Logger, config config.Config) *Coordinator {
	return &Coordinator{
		Service: &coordinatorService{},
		config:  config,
		logger:  logger,
	}
}

func (c *Coordinator) Start() {
	c.logger.Infof("Starting coordinator")

	// Set up a periodic health check for all known workers.
	go func() {
		pollTicker := time.NewTicker(workerPollFrequency)
		for _ = range pollTicker.C {
			for _, worker := range c.config.Proto().Workers {
				c.checkWorker(fmt.Sprintf("%v:%v", worker.GrpcHost, worker.GrpcPort))
			}
		}
	}()
}

func (c *Coordinator) checkWorker(address string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		c.logger.Errorf("Failed to connect to %s: %v", address, err)
	}
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	response, err := client.Status(context.TODO(), &pb_worker.StatusRequest{})
	if err != nil {
		c.logger.Errorf("Failed to retrieve status: %v", err)
	}

	health := "healthy"
	if !response.Healthy {
		health = "unhealthy"
	}
	c.logger.Infof("Worker at [%s]: [%s]", address, health)
}

type coordinatorService struct {
}

func (s *coordinatorService) Check(context.Context, *pb_coordinator.CheckRequest) (*pb_coordinator.CheckResponse, error) {
	// TODO(dino): Make up a check id, create a record for the check id.
	// TODO(dino): Look at the repository at the requested revision, work out what need to be tested.
	// TODO(dino): Pick a suitable worker (maximize caching potential), kick off the run.
	// TODO(dino): Return the check id to the caller.
	return nil, grpc.Errorf(codes.Unimplemented, "Check not implemented")
}

func (s *coordinatorService) GetStatus(context.Context, *pb_coordinator.StatusRequest) (*pb_coordinator.StatusResponse, error) {
	// TODO(dino): Lookup the requested check id
	return nil, grpc.Errorf(codes.Unimplemented, "Check not implemented")
}
