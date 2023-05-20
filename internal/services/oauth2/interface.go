package oauth2

import (
	"context"
)

type OAuth2 interface {
	MakeChallenge(ctx context.Context, challenge string) (string, error)
}
