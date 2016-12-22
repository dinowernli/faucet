package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker"

	logrus "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	workerPort = 12345

	// TODO(dino): Use flags for these.
	configFilePath = "demo/config.json"
)

var (
	// TODO(dino): Use flags for these.
	configFilePollFrequency = time.Second * 3
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	worker := worker.New()
	logger.Infof("Created worker")

	loader, err := config.NewLoader(configFilePath, configFilePollFrequency)
	if err != nil {
		logger.Fatalf("Unable to create config loader for path [%s]: %v", configFilePath, err)
	}
	coordinator := coordinator.New(logger, loader)
	logger.Infof("Created coordinator")

	go startServer(logger, worker)
	go coordinator.Start()

	// TODO(dino): Pass this into the other stuff as a shutdown channel.
	shutdown := make(chan bool)
	<-shutdown
}

func startServer(logger *logrus.Logger, worker *worker.Worker) {
	server := grpc.NewServer()
	pb_worker.RegisterWorkerServer(server, worker.Service)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", workerPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	logger.Infof("Starting worker server on port %v", workerPort)
	server.Serve(listen)
}
