package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	pb_config "dinowernli.me/faucet/proto/config"

	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(t, 1, len(config.Workers))
}

func createLoader(t *testing.T, filename string) Loader {
	loader, err := NewLoader(filename, time.Millisecond*100)
	assert.NoError(t, err, "Unable to create config loader")
	return loader
}

func tempFileWithString(t *testing.T, content string) string {
	file, err := ioutil.TempFile(os.TempDir(), "config")
	assert.NoError(t, err, "Failed to create temp file")

	err = ioutil.WriteFile(file.Name(), []byte(content), 0644)
	assert.NoError(t, err, "Failed to write to temp file")

	return file.Name()
}
