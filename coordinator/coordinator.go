package coordinator

import (
	"fmt"
	"log"

	"dinowernli.me/faucet/config"
	pb_config "dinowernli.me/faucet/proto/config"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Coordinator represents an agent in the system which implements the faucet
// coordinator service. The coordinator uses faucet workers it knows of in
// order to make sure builds get executed.
type Coordinator struct {
	configLoader config.Loader
}

// New creates a new coordinator, and is otherwise side-effect free.
func New(configLoader config.Loader) *Coordinator {
	return &Coordinator{configLoader: configLoader}
}

func (c *Coordinator) Start() {
	log.Printf("Starting coordinator")

	// TODO(dino): Make the config a field and set up atomic updates.
	var config *pb_config.Configuration
	c.configLoader.Listen(func(c *pb_config.Configuration) {
		config = c
	})

	log.Printf("Using config: %v", config)

	done := make(chan bool)
	for _, worker := range config.Workers {
		go pollStatus(fmt.Sprintf("%v:%v", worker.GrpcHost, worker.GrpcPort), done)
	}
	for _, _ = range config.Workers {
		<-done
	}
}

func pollStatus(address string, done chan (bool)) {
	log.Printf("Starting to poll address: %s", address)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", address, err)
	}

	log.Printf("Successfully dialed: %s", address)
	defer connection.Close()

	client := pb_worker.NewWorkerClient(connection)
	response, err := client.Status(context.TODO(), &pb_worker.StatusRequest{})
	if err != nil {
		log.Fatalf("Failed to retrieve status: %v", err)
	}

	log.Printf("Got response: %v", response)

	done <- true
}
