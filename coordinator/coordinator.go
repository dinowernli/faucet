package main

import (
	"fmt"
	"log"

	pb_config "dinowernli.me/faucet/proto/config"

	"google.golang.org/grpc"
)

func main() {
	log.Printf("Starting coordinator")

	config := &pb_config.Configuration{
		Workers: []*pb_config.Worker{
			&pb_config.Worker{
				GrpcHost: "localhost",
				GrpcPort: 12345,
			},
		},
	}
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

	_, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", address, err)
	} else {
		log.Printf("Successfully dialed: %s", address)
	}

	done <- true
}
