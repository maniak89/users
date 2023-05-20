package stdcrypto

import (
	"context"
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/rs/zerolog"
)

var supportedHashes = map[string]crypto.Hash{
	"NONE":                 0,
	crypto.MD5.String():    crypto.MD5,
	crypto.SHA1.String():   crypto.SHA1,
	crypto.SHA256.String(): crypto.SHA256,
	crypto.SHA512.String(): crypto.SHA512,
}

type security struct {
	hash   crypto.Hash
	config Config
}

func New(config Config) (*security, error) {
	hash, exists := supportedHashes[strings.ToUpper(config.Algo)]
	if !exists {
		supported := make([]string, 0, len(supportedHashes))
		for _, k := range supportedHashes {
			supported = append(supported, k.String())
		}
		return nil, errors.New("not supported hash: " + config.Algo + ". Supported: " + strings.Join(supported, ", "))
	}
	return &security{
		hash:   hash,
		config: config,
	}, nil
}

func (s *security) Hash(ctx context.Context, password string) (string, error) {
	if s.hash == 0 {
		return password, nil
	}
	logger := zerolog.Ctx(ctx)
	instance := s.hash.New()
	if _, err := instance.Write([]byte(password + s.config.Salt)); err != nil {
		logger.Error().Err(err).Msg("Failed write password for hash")
		return "", err
	}
	return hex.EncodeToString(instance.Sum(nil)), nil
}
