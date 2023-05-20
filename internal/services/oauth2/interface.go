package oauth2

import (
	"context"
)

type OAuth2 interface {
	MakeLogin(ctx context.Context, subject, challenge string) (string, error)
	MakeConsent(ctx context.Context, challenge string) (string, error)
}
