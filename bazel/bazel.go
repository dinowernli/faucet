package bazel

import (
	pb_workspace "dinowernli.me/faucet/proto/workspace"
)

// Client is a type which can be used to interact with Bazel.
type Client struct {
	workspace *pb_workspace.Workspace
}

func NewClient(workspace *pb_workspace.Workspace) *Client {
	return &Client{workspace: workspace}
}
