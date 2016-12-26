package config

import (
	"time"

	pb_config "dinowernli.me/faucet/proto/config"

	"github.com/Sirupsen/logrus"
)

const (
	filePollFrequency = time.Second * 3
)

// Config is a interface which provides easy access to a configuration object
// which might be changing behind the scenes.
type Config interface {
	// Proto returns an instance of the backing config proto.
	Proto() *pb_config.Configuration
}

// ForFile returns a Config instance which watches a specified file.
func ForFile(logger *logrus.Logger, filepath string) (Config, error) {
	return newLoader(logger, filepath, filePollFrequency)
}
