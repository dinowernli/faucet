package main

import (
	"fmt"
	"net"
	"os"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	workerPort = 12345

	// TODO(dino): Use flags for these.
	configFilePath = "demo/config.json"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	worker := worker.New(logger)
	logger.Infof("Created worker")

	config, err := config.ForFile(logger, configFilePath)
	if err != nil {
		logger.Fatalf("Unable to create config for path [%s]: %v", configFilePath, err)
	}
	coordinator := coordinator.New(logger, config)
	logger.Infof("Created coordinator")

	go coordinator.Start()

	// TODO(dino): Add a channel to avoid the race where the coordinator is still
	// starting by the time the first requests come in.
	go startServer(logger, worker, coordinator)

	// TODO(dino): Pass this into the other stuff as a shutdown channel.
	shutdown := make(chan bool)
	<-shutdown
}

func startServer(logger *logrus.Logger, worker *worker.Worker, coordinator *coordinator.Coordinator) {
	server := grpc.NewServer()
	pb_worker.RegisterWorkerServer(server, worker.Service)
	logger.Infof("Registered worker service")

	pb_coordinator.RegisterCoordinatorServer(server, coordinator.Service)
	logger.Infof("Registered coordinator service")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", workerPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	logger.Infof("Starting grpc server on port %v", workerPort)
	server.Serve(listen)
}
