package coordinator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sync"
	"time"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator/storage"
	pb_config "dinowernli.me/faucet/proto/config"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	pb_storage "dinowernli.me/faucet/proto/storage"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	workerPollFrequency = time.Second * 5
	checkIdSize         = 6
	healthCheckTimeout  = 100 * time.Millisecond
	executeTimeout      = 1 * time.Minute
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
	c.updateWorkerMap()
	go func() {
		// TODO(dino): rework this entirely... every worker needs a timer, etc

		pollTicker := time.NewTicker(workerPollFrequency)
		for _ = range pollTicker.C {
			c.updateWorkerMap()
		}
	}()
}

func (c *Coordinator) Check(ctx context.Context, request *pb_coordinator.CheckRequest) (*pb_coordinator.CheckResponse, error) {
	c.logger.Infof("Got check request for checkout: %v", request.Checkout)

	checkId, err := createCheckId()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to generate check id: %v", err)
	}
	c.logger.Infof("Generated check id: %s", checkId)

	workerAddress, err := c.pickWorker()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to pick a worker for check: %v", err)
	}
	c.logger.Infof("Picked worker: %v", workerAddress)

	record := &pb_storage.CheckRecord{
		Id:     checkId,
		Status: pb_storage.CheckRecord_STARTED,
	}
	err = c.storage.Put(record)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Unable to create storage record for check: %v", err)
	}
	c.logger.Infof("Created record for check")

	// TODO(dino): Make an rpc to the picked worker to kick off checking.
	execRequest := &pb_worker.ExecutionRequest{}
	c.executeCheck(workerAddress, execRequest)

	return &pb_coordinator.CheckResponse{
		CheckId: checkId,
	}, nil
}

func (s *Coordinator) GetStatus(ctx context.Context, request *pb_coordinator.StatusRequest) (*pb_coordinator.StatusResponse, error) {
	_, err := s.storage.Get(request.CheckId)
	if err != nil {
		s.logger.Errorf("Unable to load record with id [%s]: %v", request.CheckId, err)
		return nil, err
	}

	// TODO(dino): Actually use the record to populate the status response. For now, we don't
	// actually kick off any builds, so assume always pending.
	return &pb_coordinator.StatusResponse{
		Status: pb_coordinator.StatusResponse_PENDING,
	}, nil
}

// updateWorkerMap performs a health check on all known workers and udpates the
// worker pool held by this coordinator.
func (c *Coordinator) updateWorkerMap() {
	workers := c.config.Proto().Workers

	// Create a new map with an entry for every worker in the config.
	newPool := map[string]*workerStatus{}

	c.workerPoolLock.Lock()
	oldNumHealthy := numHealthy(&c.workerPool)

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
		workerStatus := c.checkWorker(workerAddress(workerProto))
		key := workerAddress(workerProto)

		existing, ok := newPool[key]
		if ok {
			existing.healthy = workerStatus.healthy
		}
	}
	newNumHealthy := numHealthy(&newPool)

	// Finally, grab the lock again and swap in the new pool.
	c.workerPoolLock.Lock()
	c.workerPool = newPool
	c.workerPoolLock.Unlock()

	if oldNumHealthy != newNumHealthy {
		c.logger.Infof("Number of healthy workers went from %d to %d", oldNumHealthy, newNumHealthy)
	}
}

// pickWorker returns the address of a healthy worker. Returns an error if we
// were unable to pick a worker.
func (c *Coordinator) pickWorker() (string, error) {
	c.workerPoolLock.Lock()
	defer c.workerPoolLock.Unlock()
	for workerAddress, workerStatus := range c.workerPool {
		if workerStatus.healthy {
			return workerAddress, nil
		}
	}
	return "", fmt.Errorf("Unable to find healthy worker")
}

func (c *Coordinator) executeCheck(address string, request *pb_worker.ExecutionRequest) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		c.logger.Errorf("Failed to connect to %s: %v", address, err)
		return
	}
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	ctx, _ := context.WithTimeout(context.Background(), executeTimeout)
	stream, err := client.Execute(ctx, request)
	if err != nil {
		c.logger.Errorf("Failed to send execute request to worker: %v", err)
		return
	}

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.logger.Errorf("Got error while reading from rpc stream: %v", err)
			return
		}

		c.logger.Infof("Got update: %v", update)
	}
}

func (c *Coordinator) checkWorker(address string) *workerStatus {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		c.logger.Errorf("Failed to connect to %s: %v", address, err)
		return &workerStatus{healthy: false}
	}
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	ctx, _ := context.WithTimeout(context.Background(), healthCheckTimeout)
	response, err := client.Status(ctx, &pb_worker.StatusRequest{})
	if err != nil {
		c.logger.Errorf("Failed to retrieve status: %v", err)
		return &workerStatus{healthy: false}
	}

	health := "healthy"
	if !response.Healthy {
		health = "unhealthy"
	}
	c.logger.Debugf("Worker at [%s]: [%s]", address, health)
	return &workerStatus{healthy: response.Healthy}
}

func workerAddress(proto *pb_config.Worker) string {
	return fmt.Sprintf("%v:%v", proto.GrpcHost, proto.GrpcPort)
}

func numHealthy(workerPool *map[string]*workerStatus) int {
	result := 0
	for _, status := range *workerPool {
		if status.healthy {
			result++
		}
	}
	return result
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
