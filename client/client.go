package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

	rpcTimeout = time.Second * 3
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

	repoMeta, err := getRepoMetadata(repoRoot)
	if err != nil {
		logger.Errorf("Unable to get repository metadata: %v", err)
		return
	}

	cloneUrl := repoMeta.GitRepo.CloneUrl
	if cloneUrl == "" {
		logger.Errorf("Got invalid clone url: %s", cloneUrl)
		return
	}

	commitHash, err := getCommitHash()
	if err != nil {
		logger.Errorf("Unable to get commit hash: %v", err)
		return
	}

	checkRequest := &pb_coordinator.CheckRequest{
		Checkout: &pb_workspace.Checkout{
			Revision:   revisionProto(commitHash),
			Repository: repoProto(cloneUrl),
		},
	}
	logger.Infof("Sending CheckRequest: %v", checkRequest)

	// TODO(dino): Add SSL and context deadlines.
	connection, err := grpc.Dial(*flagCoordinator, grpc.WithInsecure(), grpc.WithTimeout(rpcTimeout))
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

func revisionProto(gitCommitHash string) *pb_workspace.Revision {
	return &pb_workspace.Revision{
		&pb_workspace.Revision_GitRevision_{
			GitRevision: &pb_workspace.Revision_GitRevision{
				CommitHash: gitCommitHash,
			},
		},
	}
}

func repoProto(cloneUrl string) *pb_workspace.Repository {
	return &pb_workspace.Repository{
		&pb_workspace.Repository_GitRepository_{
			&pb_workspace.Repository_GitRepository{
				CloneUrl: cloneUrl,
			},
		},
	}
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
