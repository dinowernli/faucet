package checkout

import (
	"fmt"

	pb_workspace "dinowernli.me/faucet/proto/workspace"
)

// Checkout represents a directory tree in a specific requested state. There
// are acquired recources associated with a checkout instance (e.g., a backing
// version control repository), so it is the user's responsibility to dispose
// of the checkout instance after they are done using it.
type Checkout struct {
	RootPath string
}

// Close returns underlying resources associated with this instance. After
// calling this method, the instance must no longer be used.
func (c *Checkout) Close() {
}

// CheckoutProvider is an interface through which Checkouts can be obtained.
type CheckoutProvider interface {
	Get(proto *pb_workspace.Checkout) (*Checkout, error)
}

func NewProvider() CheckoutProvider {
	return &checkoutProvider{}
}

type checkoutProvider struct {
}

func (c *checkoutProvider) Get(proto *pb_workspace.Checkout) (*Checkout, error) {
	// TODO(dino): Allocate a temporary directory, fetch the repo and the commit.
	return nil, fmt.Errorf("checkoutProvider.Get not implemented")
}
