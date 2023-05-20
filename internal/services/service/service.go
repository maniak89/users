package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"users/internal/crypto"
	"users/internal/services/oauth2"
	"users/internal/storage"
)

var ErrNoUser = errors.New("user not found or password mismatch")

type Service struct {
	Storage storage.Storage
	Crypto  crypto.Crypto
	OAuth2  oauth2.OAuth2
}

func (s *Service) Login(ctx context.Context, login, password, challenge string) (string, error) {
	logger := log.Ctx(ctx).With().Str("login", login).Logger()
	ctx = logger.WithContext(ctx)
	user, err := s.Storage.Find(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			logger.Debug().Err(err).Msg("User not found")
			return "", ErrNoUser
		}
		logger.Error().Err(err).Msg("Failed find user")
		return "", err
	}
	hash, err := s.Crypto.Hash(ctx, password)
	if err != nil {
		logger.Error().Err(err).Msg("Failed calc hash")
		return "", err
	}
	if user.Password != hash {
		logger.Debug().Msg("password mismatch")
		return "", ErrNoUser
	}
	redirectTo, err := s.OAuth2.MakeLogin(ctx, user.ID, challenge)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make login request")
		return "", err
	}
	if err := s.Storage.UpdateLastLogin(ctx, login); err != nil {
		logger.Warn().Err(err).Msg("Failed update last login")
	}
	return redirectTo, nil
}

func (s *Service) Consent(ctx context.Context, challenge string) (string, error) {
	logger := log.Ctx(ctx).With().Str("challenge", challenge).Logger()
	ctx = logger.WithContext(ctx)
	redirectTo, err := s.OAuth2.MakeConsent(ctx, challenge)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make challenge request")
		return "", err
	}
	return redirectTo, nil
}
