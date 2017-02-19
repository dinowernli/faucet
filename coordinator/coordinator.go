package coordinator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator/storage"
	pb_config "dinowernli.me/faucet/proto/config"
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
	config         config.Config
	logger         *logrus.Logger
	storage        storage.Storage
	workerPool     map[string]*workerStatus
	workerPoolLock *sync.Mutex
}

// workerStatus is the type used by the coordinator to keep track of a worker.
// The coordinator maintain a status object for every worker in the config, but
// might decide not to send any traffic to a worker because it is unhealthy.
type workerStatus struct {
	healthy bool
}

// New creates a new coordinator, and is otherwise side-effect free.
func New(logger *logrus.Logger, config config.Config) *Coordinator {
	return &Coordinator{
		config:         config,
		logger:         logger,
		storage:        storage.NewInMemory(),
		workerPool:     map[string]*workerStatus{},
		workerPoolLock: &sync.Mutex{},
	}
}

func (c *Coordinator) Start() {
	c.logger.Infof("Starting coordinator")

	// Set up a periodic health check for all known workers.
	go func() {
		// TODO(dino): rework this entirely... every worker needs a timer, etc

		pollTicker := time.NewTicker(workerPollFrequency)
		for _ = range pollTicker.C {
			workers := c.config.Proto().Workers

			// Create a new map with an entry for every worker in the config.
			newPool := map[string]*workerStatus{}
			c.workerPoolLock.Lock()
			for _, workerProto := range workers {
				key := workerAddress(workerProto)
				existing, ok := c.workerPool[key]
				if ok {
					newPool[key] = existing
				} else {
					newPool[key] = &workerStatus{healthy: false}
				}
			}
			c.workerPoolLock.Unlock()

			// Then, update worker statuses in-place based on health checks.
			for _, workerProto := range workers {
				workerStatus := c.checkWorker(workerProto)
				key := workerAddress(workerProto)

				existing, ok := c.workerPool[key]
				if ok {
					existing.healthy = workerStatus.healthy
				}
			}

			// Finally, grab the lock again and swap in the new pool.
			c.workerPoolLock.Lock()
			c.workerPool = newPool
			c.workerPoolLock.Unlock()
		}
	}()
}

func (s *Coordinator) Check(ctx context.Context, request *pb_coordinator.CheckRequest) (*pb_coordinator.CheckResponse, error) {
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

func (s *Coordinator) GetStatus(ctx context.Context, request *pb_coordinator.StatusRequest) (*pb_coordinator.StatusResponse, error) {
	_, err := s.storage.Get(request.CheckId)
	if err != nil {
		s.logger.Errorf("Unable to load record with id [%s]: %v", request.CheckId, err)
		return nil, err
	}

	// TODO(dino): Actually use the record to populate the status response.
	return &pb_coordinator.StatusResponse{}, nil
}

func (c *Coordinator) checkWorker(proto *pb_config.Worker) *workerStatus {
	address := workerAddress(proto)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		c.logger.Errorf("Failed to connect to %s: %v", address, err)
		return &workerStatus{healthy: false}
	}
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	// TODO(dino): Set a deadline here.
	response, err := client.Status(context.TODO(), &pb_worker.StatusRequest{})
	if err != nil {
		c.logger.Errorf("Failed to retrieve status: %v", err)
		return &workerStatus{healthy: false}
	}

	health := "healthy"
	if !response.Healthy {
		health = "unhealthy"
	}
	c.logger.Infof("Worker at [%s]: [%s]", address, health)
	return &workerStatus{healthy: response.Healthy}
}

func workerAddress(proto *pb_config.Worker) string {
	return fmt.Sprintf("%v:%v", proto.GrpcHost, proto.GrpcPort)
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
