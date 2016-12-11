package main

import (
	"log"
	"net"

	"dinowernli.me/faucet/demo"
	pb_config "dinowernli.me/faucet/proto/config"
	pb_worker "dinowernli.me/faucet/proto/service/worker"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type workerService struct {
}

func (s *workerService) Status(context context.Context, request *pb_worker.StatusRequest) (*pb_worker.StatusResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Not implemented")
}

func main() {
	log.Printf(demo.Foo())

	_ = &pb_config.Configuration{}
	_ = &pb_worker.StatusRequest{}

	server := grpc.NewServer()
	pb_worker.RegisterWorkerServer(server, &workerService{})

	listen, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting worker server on port 12345")
	server.Serve(listen)
}
