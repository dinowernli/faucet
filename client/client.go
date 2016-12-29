package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	pb_client "dinowernli.me/faucet/proto/client"
	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"
	pb_workspace "dinowernli.me/faucet/proto/workspace"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/jsonpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// TODO(dino): Check for the existence of the binaries, allow overriding.
	gitBinary            = "git"
	repoMetadataFilename = "FAUCET"
)

var (
	flagCoordinator = flag.String("coordinator", "localhost:12345", "the grpc endpoint of the coordinator")
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	logger.Infof("Using coordinator: %s", *flagCoordinator)

	// TODO(dino): Check that there are no uncommited changes.
	// TODO(dino): Have better messaging around common failure cases

	repoRoot, err := getRepoRoot()
	if err != nil {
		logger.Errorf("Unable to get repository root: %v", err)
		return
	}
	logger.Infof("Using repository root: %s", repoRoot)

	repoMeta, err := getRepoMetadata(repoRoot)
	if err != nil {
		logger.Errorf("Unable to get repository metadata: %v", err)
		return
	}
	logger.Infof("Using repository metadata: %s", repoMeta)

	commitHash, err := getCommitHash()
	if err != nil {
		logger.Errorf("Unable to get commit hash: %v", err)
		return
	}
	logger.Infof("Using commit hash: %s", commitHash)
	revision := &pb_workspace.Revision{
		&pb_workspace.Revision_GitRevision_{
			GitRevision: &pb_workspace.Revision_GitRevision{
				CommitHash: commitHash,
			},
		},
	}

	request := &pb_coordinator.CheckRequest{
		Checkout: &pb_workspace.Checkout{
			Revision: revision,
		},
	}
	logger.Infof("Request: %v", request)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	connection, err := grpc.Dial(*flagCoordinator, opts...)
	if err != nil {
		logger.Errorf("Failed to connect to %s: %v", *flagCoordinator, err)
		return
	}
	defer connection.Close()

	client := pb_coordinator.NewCoordinatorClient(connection)
	response, err := client.Check(context.TODO(), &pb_coordinator.CheckRequest{})
	if err != nil {
		logger.Errorf("Failed to retrieve status: %v", err)
		return
	}

	logger.Infof("Got response: %v", response)
}

func getRepoMetadata(repoRoot string) (*pb_client.RepositoryMetadata, error) {
	filename := filepath.Join(repoRoot, repoMetadataFilename)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file %s: %v", filename, err)
	}

	result := &pb_client.RepositoryMetadata{}
	err = jsonpb.UnmarshalString(string(fileBytes), result)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse proto file: %v", err)
	}

	return result, nil
}

func getRepoRoot() (string, error) {
	outBytes, err := exec.Command(gitBinary, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("Unable to execute git: %v", err)
	}
	return strings.Trim(string(outBytes), "\r\n"), nil
}

func getCommitHash() (string, error) {
	outBytes, err := exec.Command(gitBinary, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("Unable to execute git: %v", err)
	}
	return strings.Trim(string(outBytes), "\r\n"), nil
}
