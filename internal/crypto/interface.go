package crypto

import (
	"context"
)

type Crypto interface {
	Hash(ctx context.Context, password string) (string, error)
}
