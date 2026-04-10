package accountschecker

import (
	"context"
)

// Checker defines the interface for the accounts history checker.
type Checker interface {
	Run(ctx context.Context) error
}

// APIClient defines the interface for communicating with the deep history API.
type APIClient interface {
	GetHighestNonce(ctx context.Context) (uint64, error)
	GetAccountData(ctx context.Context, address string, nonce uint64) (*AccountData, error)
}
