package hydra

import (
	"context"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/rs/zerolog/log"
)

type service struct {
	client *hydra.APIClient
}

func New(config Config) *service {
	conf := hydra.NewConfiguration()
	conf.Host = config.HydraAdminUrl
	return &service{
		client: hydra.NewAPIClient(conf),
	}
}

func (s *service) MakeChallenge(ctx context.Context, challenge string) (string, error) {
	logger := log.Ctx(ctx)
	request := s.client.OAuth2Api.AcceptOAuth2ConsentRequest(ctx)
	redirectTo, _, err := request.ConsentChallenge(challenge).Execute()
	if err != nil {
		logger.Error().Err(err).Msg("Failed make challenge request")
		return "", err
	}
	return redirectTo.RedirectTo, nil
}
