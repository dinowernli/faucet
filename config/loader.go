package config

import (
	pb_config "dinowernli.me/faucet/proto/config"
)

// Loader tracks a config proto from a single source (e.g., a file) and
// dispatches notifications when the config changes.
type Loader interface {
	// Listen attaches a listener to this config loader. This call immediately
	// triggers a callback for the current config.
	Listen(func(*pb_config.Configuration))
}

// NewLoader creates a loader which watches a config file.
func NewLoader(filepath string) Loader {
	// TODO(dino): Open the file, load config, install a watching goroutine.
	return &loader{}
}

type loader struct {
}

func (l *loader) Listen(callback func(*pb_config.Configuration)) {
}
