package bazel

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// TODO(dino): Figure out what to do about workspaces.

// Client is a type which can be used to interact with Bazel.
type Client interface {
	// Run executes a bazel build.
	Run(targets []string) error
}

type client struct {
	logger *logrus.Logger
}

func NewClient(logger *logrus.Logger) Client {
	return &client{logger: logger}
}

func (c *client) Run(targets []string) error {
	// Dummy implementation for now.
	c.logger.Infof("Fake building with Bazel...", targets)
	time.Sleep(time.Second * 3)
	c.logger.Infof("Finished fake building")
	return nil
}
