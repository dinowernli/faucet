package config

import (
	"io/ioutil"
	"sync"
	"time"

	pb_config "dinowernli.me/faucet/proto/config"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type callback func(*pb_config.Configuration)

// Loader tracks a config proto from a single source (e.g., a file) and
// dispatches notifications when the config changes.
type Loader interface {
	// Listen attaches a listener to this config loader. This call immediately
	// triggers a callback for the current config. A registered callback will not
	// get executed concurrently with itself, but might get executed on different
	// goroutines over time.
	Listen(callback)
}

// newLoader creates a loader which watches a config file.
func newLoader(logger *logrus.Logger, filepath string, pollFrequency time.Duration) (Loader, error) {
	initialConfig, err := readFile(filepath)
	if err != nil {
		return nil, err
	}

	result := &loader{
		logger:        logger,
		pollFrequency: pollFrequency,
		config:        initialConfig,
		configLock:    &sync.Mutex{},
		callbacks:     []callback{},
		callbacksLock: &sync.Mutex{},
	}

	// Kick off a goroutin to monitor the file.
	go result.pollFile(filepath)

	return result, nil
}

type loader struct {
	logger        *logrus.Logger
	pollFrequency time.Duration
	config        *pb_config.Configuration
	configLock    *sync.Mutex
	callbacks     []callback
	callbacksLock *sync.Mutex
}

func (l *loader) Listen(cb callback) {
	cb(l.currentConfig())

	// Only add the callback to the list once it's done executing to make sure
	// it doesn't get executed concurrently by the polling goroutine.
	l.callbacksLock.Lock()
	defer l.callbacksLock.Unlock()
	l.callbacks = append(l.callbacks, cb)
}

func (l *loader) currentConfig() *pb_config.Configuration {
	l.configLock.Lock()
	defer l.configLock.Unlock()
	return l.config
}

// updateConfig sets the current config. Returns true if the new config is
// different from the one set previously.
func (l *loader) updateConfig(config *pb_config.Configuration) bool {
	l.configLock.Lock()
	defer l.configLock.Unlock()

	result := !proto.Equal(config, l.config)
	l.config = config
	return result
}

func (l *loader) callbacksSnapshot() []callback {
	l.callbacksLock.Lock()
	defer l.callbacksLock.Unlock()

	result := []callback{}
	for _, cb := range l.callbacks {
		result = append(result, cb)
	}
	return result
}

func (l *loader) pollFile(filepath string) {
	ticker := time.NewTicker(l.pollFrequency)
	for _ = range ticker.C {
		config, err := readFile(filepath)
		if err != nil {
			l.logger.Warnf("Polling file [%s] failed: %v", filepath, err)
			continue
		}

		if l.updateConfig(config) {
			l.logger.Infof("Updated config")
			for _, cb := range l.callbacksSnapshot() {
				cb(config)
			}
		}
	}
}

func readFile(filepath string) (*pb_config.Configuration, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	result := &pb_config.Configuration{}
	err = jsonpb.UnmarshalString(string(bytes), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
