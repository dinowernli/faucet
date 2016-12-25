package main

import (
	"os"

	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"

	"github.com/Sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	_ = &pb_coordinator.CheckRequest{}
	// TODO(dino): Take a few command line args to populate the request body as
	// well as to find the coordinator, then send request to coordinator.

	logger.Infof("Hello world")
}
