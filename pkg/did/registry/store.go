package registry

import (
	"context"
)

// Store defines persistence operations for DID documents.
type Store interface {
	Put(ctx context.Context, did string, document []byte) error
	Get(ctx context.Context, did string) ([]byte, error)
	Has(ctx context.Context, did string) (bool, error)
}
