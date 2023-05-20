package stdcrypto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_security_Hash(t *testing.T) {
	instr := "test string"
	tests := []struct {
		config  Config
		want    string
		wantErr bool
	}{
		{
			config: Config{Algo: "none"},
			want:   instr,
		},
		{
			config: Config{Algo: "md5"},
			want:   "6f8db599de986fab7a21625b7916589c",
		},
		{
			config: Config{Algo: "sha-1"},
			want:   "661295c9cbf9d6b2f6428414504a8deed3020641",
		},
		{
			config: Config{Algo: "sha-256"},
			want:   "d5579c46dfcc7f18207013e65b44e4cb4e2c2298f4ac457ba8f82743f31e930b",
		},
		{
			config: Config{Algo: "sha-512"},
			want:   "10e6d647af44624442f388c2c14a787ff8b17e6165b83d767ec047768d8cbcb71a1a3226e7cc7816bc79c0427d94a9da688c41a3992c7bf5e4d7cc3e0be5dbac",
		},
	}
	for _, tt := range tests {
		t.Run(tt.config.Algo, func(t *testing.T) {
			s, err := New(tt.config)
			assert.NoError(t, err)
			got, err := s.Hash(context.Background(), instr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
