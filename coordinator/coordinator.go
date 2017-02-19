package coordinator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator/storage"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	workerPollFrequency = time.Second * 2
	checkIdSize         = 6
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
		Service: &coordinatorService{
			logger:  logger,
			storage: storage.NewInMemory(),
		},
		config: config,
		logger: logger,
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
		return
	}
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	response, err := client.Status(context.TODO(), &pb_worker.StatusRequest{})
	if err != nil {
		c.logger.Errorf("Failed to retrieve status: %v", err)
		return
	}

	health := "healthy"
	if !response.Healthy {
		health = "unhealthy"
	}
	c.logger.Infof("Worker at [%s]: [%s]", address, health)
}

type coordinatorService struct {
	logger  *logrus.Logger
	storage storage.Storage
}

func (s *coordinatorService) Check(ctx context.Context, request *pb_coordinator.CheckRequest) (*pb_coordinator.CheckResponse, error) {
	s.logger.Infof("Got check request: %v", request)

	checkId, err := createCheckId()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to generate check id: %v", err)
	}
	s.logger.Infof("Generated check id: %s", checkId)

	// TODO(dino): Pick a suitable worker (maximize caching potential), kick off the run.
	// TODO(dino): Return the check id to the caller.
	return nil, grpc.Errorf(codes.Unimplemented, "Check not implemented")
}

func (s *coordinatorService) GetStatus(ctx context.Context, request *pb_coordinator.StatusRequest) (*pb_coordinator.StatusResponse, error) {
	_, err := s.storage.Get(request.CheckId)
	if err != nil {
		s.logger.Errorf("Unable to load record with id [%s]: %v", request.CheckId, err)
		return nil, err
	}

	// TODO(dino): Actually use the record to populate the status response.
	return &pb_coordinator.StatusResponse{}, nil
}

func createCheckId() (string, error) {
	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Unable to read random bytes into buffer: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buffer)
	return encoded[0 : checkIdSize-1], nil
}
