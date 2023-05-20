package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"users/internal/services/rest/handlers"
	serviceImpl "users/internal/services/service"
)

type service struct {
	config Config
	srv    *http.Server
}

const xRequestID = "X-Request-Id"

func New(ctx context.Context, config Config, log zerolog.Logger, serviceImpl *serviceImpl.Service) (*service, error) {
	r := chi.NewRouter()
	r.Use(
		hlog.NewHandler(log),
		hlog.MethodHandler("method"),
		hlog.URLHandler("url"),
		hlog.RequestIDHandler("x_request_id", xRequestID),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			zerolog.Ctx(ctx).Trace().Str("method", r.Method).Str("url", r.URL.String()).Str("x_request_id", r.Header.Get(xRequestID)).Int("status", status).Int("size", size).Dur("duration", duration).Msg("request processed")
		}),
		middleware.Recoverer,
	)
	service := service{
		config: config,
		srv:    &http.Server{Addr: config.Address, Handler: r},
	}

	r.Head("/", service.Health)
	r.Route("/login", func(r chi.Router) {
		r.Get("/", handlers.LoginGet)
		r.Post("/", handlers.LoginPost(serviceImpl))
	})
	r.Get("/consent", handlers.ConsentGet(serviceImpl))

	return &service, nil
}

func (s *service) Run(ctx context.Context, ready func()) error {
	logger := log.Ctx(ctx)
	logger.Info().Str("address", s.srv.Addr).Msg("Start listening")
	defer func() {
		logger.Info().Msg("Stop listening")
	}()
	ready()
	if err := s.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		logger.Error().Err(err).Msg("Failed start listening")
		return err
	}

	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	logger := log.Ctx(ctx)

	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed shutdown")
		return err
	}

	return nil
}
