package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"dinowernli.me/faucet/config"
	"dinowernli.me/faucet/coordinator"
	pb_worker "dinowernli.me/faucet/proto/service/worker"
	"dinowernli.me/faucet/worker"

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
	worker := worker.New()
	log.Printf("Created worker")

	loader, err := config.NewLoader(configFilePath, configFilePollFrequency)
	if err != nil {
		log.Fatalf("Unable to create config loader for path [%s]: %v", configFilePath, err)
	}
	coordinator := coordinator.New(loader)
	log.Printf("Created coordinator")

	go startServer(worker)
	go coordinator.Start()

	// TODO(dino): Pass this into the other stuff as a shutdown channel.
	shutdown := make(chan bool)
	<-shutdown
}

func startServer(worker *worker.Worker) {
	server := grpc.NewServer()
	pb_worker.RegisterWorkerServer(server, worker.Service)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", workerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting worker server on port %v", workerPort)
	server.Serve(listen)
}
