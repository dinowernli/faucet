package repository

import (
	"fmt"
	"os/exec"

	pb_workspace "dinowernli.me/faucet/proto/workspace"

	"github.com/Sirupsen/logrus"
)

const (
	gitBinaryPath = "git"
)

// Client is a type which can be used to interact with a repository.
type Client interface {
	// Checkout takes the description of a version control checkout and returns
	// the (absolute) path to a directory with the corresponding files.
	Checkout(*pb_workspace.Checkout) (string, error)
}

type client struct {
	logger *logrus.Logger
}

func NewClient(logger *logrus.Logger) Client {
	return &client{logger: logger}
}

func (c *client) Checkout(checkout *pb_workspace.Checkout) (string, error) {
	destination, err := c.clone(checkout.Repository)
	if err != nil {
		return "", fmt.Errorf("Unable to clone repo %v: %v", checkout.Repository, err)
	}

	// TODO(dino): CD in and checkout the commit hash.
	return destination, nil
}

// clone take the supplied repo description and produces a directory with the
// repo checked out.
func (c *client) clone(repoMeta *pb_workspace.Repository) (string, error) {
	gitRepoMeta := repoMeta.GetGitRepository()
	if gitRepoMeta == nil {
		return "", fmt.Errorf("Expected git repository, but got: %v", repoMeta.GetRepo())
	}

	// TODO(dino): Create an actual temporary directory.
	destination := "/tmp/yolo"

	// TODO(dino): Check that this is really how git works (?!)
	url := gitRepoMeta.CloneUrl
	err := c.logAndRun(gitBinaryPath, "clone", url, destination)
	if err != nil {
		return "", fmt.Errorf("Unable to execute git clone [%s]: %v", url, err)
	}

	return destination, nil
}

func (c *client) logAndRun(name string, arg ...string) error {
	c.logger.Infof("Running command: [%s]", "dsf")
	return exec.Command(name, arg...).Run()
}
