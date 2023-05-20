package hydra

import (
	"context"
	"net/url"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/rs/zerolog/log"
)

type service struct {
	client *hydra.APIClient
}

func New(config Config) (*service, error) {
	conf := hydra.NewConfiguration()
	uri, err := url.Parse(config.HydraAdminUrl)
	if err != nil {
		return nil, err
	}
	conf.Host = uri.Host
	conf.Scheme = uri.Scheme
	return &service{
		client: hydra.NewAPIClient(conf),
	}, nil
}

func (s *service) MakeChallenge(ctx context.Context, challenge string) (string, error) {
	logger := log.Ctx(ctx)
	request := s.client.OAuth2Api.AcceptOAuth2LoginRequest(ctx)
	redirectTo, _, err := request.LoginChallenge(challenge).Execute()
	if err != nil {
		logger.Error().Err(err).Msg("Failed make challenge request")
		return "", err
	}
	return redirectTo.RedirectTo, nil
}
