package main

import (
	"os"

	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"

	"github.com/Sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	_ = &pb_coordinator.ChangeValidationRequest{}

	logger.Infof("Hello world")
}
