package checkout

import (
	"testing"

	pb_workspace "dinowernli.me/faucet/proto/workspace"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	provider := NewProvider()
	_, err := provider.Get(&pb_workspace.Checkout{})
	assert.Error(t, err)
}
