package coordinator

import (
	"fmt"

	"dinowernli.me/faucet/config"
	pb_config "dinowernli.me/faucet/proto/config"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Coordinator represents an agent in the system which implements the faucet
// coordinator service. The coordinator uses faucet workers it knows of in
// order to make sure builds get executed.
type Coordinator struct {
	Service      *coordinatorService
	logger       *logrus.Logger
	configLoader config.Loader
}

// New creates a new coordinator, and is otherwise side-effect free.
func New(logger *logrus.Logger, configLoader config.Loader) *Coordinator {
	return &Coordinator{
		Service:      &coordinatorService{},
		logger:       logger,
		configLoader: configLoader,
	}
}

func (c *Coordinator) Start() {
	c.logger.Infof("Starting coordinator")

	// TODO(dino): Make the config a field and set up atomic updates.
	var config *pb_config.Configuration
	c.configLoader.Listen(func(c *pb_config.Configuration) {
		config = c
	})

	c.logger.Infof("Using config: %v", config)

	done := make(chan bool)
	for _, worker := range config.Workers {
		go c.pollStatus(fmt.Sprintf("%v:%v", worker.GrpcHost, worker.GrpcPort), done)
	}
	for _, _ = range config.Workers {
		<-done
	}
}

func (c *Coordinator) pollStatus(address string, done chan (bool)) {
	c.logger.Infof("Starting to poll address: %s", address)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		c.logger.Fatalf("Failed to connect to %s: %v", address, err)
	}

	c.logger.Infof("Successfully dialed: %s", address)
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	response, err := client.Status(context.TODO(), &pb_worker.StatusRequest{})
	if err != nil {
		c.logger.Fatalf("Failed to retrieve status: %v", err)
	}

	c.logger.Infof("Got response: %v", response)

	done <- true
}

type coordinatorService struct {
}

func (s *coordinatorService) ValidateChange(context.Context, *pb_coordinator.ChangeValidationRequest) (*pb_coordinator.ChangeValidationResponse, error) {
	return nil, nil
}

func (s *coordinatorService) GetValidationStatus(context.Context, *pb_coordinator.ValidationStatusRequest) (*pb_coordinator.ValidationStatusResponse, error) {
	return nil, nil
}
