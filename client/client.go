package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	pb_coordinator "dinowernli.me/faucet/proto/service/coordinator"

	"github.com/Sirupsen/logrus"
)

const (
	bazelBinary = "bazel"
	gitBinary   = "git"
)

var (
	workspaceNameRegex = regexp.MustCompile("workspace\\(name = \"([a-z_]+)\"\\)")
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stderr

	// TODO(dino): Check that there are no uncommited changes.

	commitHash, err := getCommitHash()
	if err != nil {
		logger.Errorf("Unable to get commit hash: %v", err)
		return
	}

	workspaceName, err := getWorkspaceName()
	if err != nil {
		logger.Errorf("Unable to get workspace name: %v", err)
		return
	}

	request := &pb_coordinator.CheckRequest{
		WorkspaceName: workspaceName,
		GitCommit:     commitHash,
	}
	logger.Infof("Request: %v", request)
}

func getWorkspaceName() (string, error) {
	// Find the workspace root.
	workspacePath, err := getWorkspacePath()
	if err != nil {
		return "", fmt.Errorf("Unable to get workspace path: %v", err)
	}

	// Read the WORKSPACE file.
	workspaceFilePath := filepath.Join(workspacePath, "WORKSPACE")
	fileBytes, err := ioutil.ReadFile(workspaceFilePath)
	if err != nil {
		return "", fmt.Errorf("Unable to read WORKSPACE file [%s]: %v", workspaceFilePath, err)
	}

	// Find the workspace name inside it.
	results := workspaceNameRegex.FindStringSubmatch(string(fileBytes))

	// Expect the result to contain exactly two entries: one for the entire match and one for
	// the the actual workspace name.
	if len(results) != 2 {
		return "", fmt.Errorf("Expected 2 matches, but got %d. Matches: %q", len(results), results)
	}
	return results[1], nil
}

func getWorkspacePath() (string, error) {
	outBytes, err := exec.Command(bazelBinary, "info", "workspace").Output()
	if err != nil {
		return "", fmt.Errorf("Unable to get workspace path: %v", err)
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
