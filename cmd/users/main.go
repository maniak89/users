package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"

	"users/internal/crypto/stdcrypto"
	"users/internal/log"
	"users/internal/services"
	"users/internal/services/oauth2/hydra"
	"users/internal/services/rest"
	serviceImpl "users/internal/services/service"
	"users/internal/storage/sql"
)

type config struct {
	Logger  log.Config
	Storage sql.Config
	Crypto  stdcrypto.Config
	Rest    rest.Config
	OAuth2  hydra.Config
}

const signalChLen = 10

func main() {

	var cfg config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot decode config envs")
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot init logger")
	}

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))

	g := &run.Group{}
	{
		stop := make(chan os.Signal, signalChLen)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}

	crypto, err := stdcrypto.New(cfg.Crypto)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot init crypto")
	}
	oauth2 := hydra.New(cfg.OAuth2)
	storage := sql.New(cfg.Storage)
	serviceImpl := &serviceImpl.Service{
		Storage: storage,
		OAuth2:  oauth2,
		Crypto:  crypto,
	}

	orderRunner := services.OrderRunner{}
	restService, err := rest.New(ctx, cfg.Rest, logger.With().Str("role", "rest").Logger(), serviceImpl)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create rest service")
	}
	if err := orderRunner.SetupService(ctx, restService, "rest", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup rest service")
	}

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")

}
