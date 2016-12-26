package config

import (
	"fmt"
	"sync"
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
	loader, err := newLoader(logger, filepath, filePollFrequency)
	if err != nil {
		return nil, fmt.Errorf("Unable to create config loader: %v", err)
	}
	return forLoader(loader), nil
}

// forLoader returns a new Config instance backed by the supplied config loader.
func forLoader(loader Loader) Config {
	result := &config{
		protoLock: &sync.Mutex{},
		proto:     nil, // Initialized immediately below by listening.
	}
	loader.Listen(func(proto *pb_config.Configuration) {
		result.accept(proto)
	})
	return result
}

type config struct {
	proto     *pb_config.Configuration
	protoLock *sync.Mutex
}

func (c *config) Proto() *pb_config.Configuration {
	// TODO(dino): Use a RW-lock here.
	c.protoLock.Lock()
	defer c.protoLock.Unlock()
	return c.proto
}

func (c *config) accept(proto *pb_config.Configuration) {
	c.protoLock.Lock()
	defer c.protoLock.Unlock()
	c.proto = proto
}
