package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	pb_config "dinowernli.me/faucet/proto/config"
)

// TODO(dino): Get the stretchr/testify testing lib and use proper assertions.

func TestLoader_Creation(t *testing.T) {
	filename := tempFileWithString(t, "{}")
	defer os.Remove(filename)

	createLoader(t, filename)
}

func TestLoader_InitialCallback(t *testing.T) {
	filename := tempFileWithString(t, "{ \"workers\": [ {} ] }")
	defer os.Remove(filename)

	loader := createLoader(t, filename)
	var config *pb_config.Configuration
	loader.Listen(func(c *pb_config.Configuration) {
		config = c
	})

	if len(config.Workers) != 1 {
		t.Errorf("Unexpected number of workers")
	}
}

func createLoader(t *testing.T, filename string) Loader {
	loader, err := NewLoader(filename, time.Millisecond*100)
	if err != nil {
		t.Errorf("Unable to create loader: %v", err)
	}
	return loader
}

func tempFileWithString(t *testing.T, content string) string {
	file, err := ioutil.TempFile(os.TempDir(), "config")
	if err != nil {
		t.Errorf("Failed to create file: %v", err)
	}

	err = ioutil.WriteFile(file.Name(), []byte(content), 0644)
	if err != nil {
		t.Errorf("Failed ti write to file: %v", err)
	}

	return file.Name()
}
