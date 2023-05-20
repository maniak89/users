package sql

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects"

	storageModels "users/internal/models/storage"
	storagePkg "users/internal/storage"
)

type storage struct {
	config     Config
	connection *sql.DB
	db         *reform.Querier
	reformDB   reform.DBTXContext
	driver     string
}

func New(config Config) *storage {
	return &storage{
		config: config,
	}
}

func (s *storage) Connect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	parsedConnectionString, err := url.Parse(s.config.ConnectionString)
	if err != nil {
		logger.Error().Err(err).Str("connection_string", s.config.ConnectionString).Msg("Failed parse connection string")
		return err
	}
	s.driver = parsedConnectionString.Scheme
	var found bool
	for _, driver := range sql.Drivers() {
		if driver == s.driver {
			found = true
			break
		}
	}
	if !found {
		err := errors.New("not supported db driver: " + s.driver)
		logger.Error().Err(err).Msg("")
		return err
	}
	connectionString := s.config.ConnectionString
	// sqlite3 не работает, если в начале connectionString - схема.
	// а без нее не работает postgres(не парсит аргументы)
	if s.driver == "sqlite3" {
		connectionString = s.config.ConnectionString[len(s.driver)+3:]
	}
	sqlDB, err := sql.Open(s.driver, connectionString)
	if err != nil {
		logger.Error().Err(err).Msg("Failed open sql connection")
		return err
	}
	s.connection = sqlDB

	t := reform.NewDB(sqlDB, dialects.ForDriver(s.driver), reform.NewPrintfLogger(logger.Printf))
	s.db = t.Querier
	s.reformDB = t
	return nil
}

func (s *storage) Disconnect(ctx context.Context) error {
	logger := log.Ctx(ctx)
	if s.connection == nil {
		err := storagePkg.ErrInvalidState
		logger.Error().Err(err).Msg("Invalid state of connection")
		return err
	}
	if err := s.connection.Close(); err != nil {
		logger.Error().Err(err).Msg("Failed close of connection")
		return err
	}
	s.connection = nil
	return nil
}

func (s *storage) Find(ctx context.Context, login string) (*storageModels.User, error) {
	logger := log.Ctx(ctx)
	var result storageModels.User
	if err := s.db.WithContext(ctx).SelectOneTo(&result, "WHERE login = "+s.db.Placeholder(1), login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug().Msg("User not found")
			return nil, storagePkg.ErrNoRows
		}
		logger.Error().Err(err).Msg("Failed find user")
		return nil, err
	}
	return &result, nil
}

func (s *storage) UpdateLastLogin(ctx context.Context, login string) error {
	logger := log.Ctx(ctx)
	if _, err := s.db.WithContext(ctx).Exec("update users set last_login=now() where login = "+s.db.Placeholder(1), login); err != nil {
		logger.Error().Err(err).Msg("Failed update user")
		return err
	}
	return nil
}
