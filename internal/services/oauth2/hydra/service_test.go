package hydra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	client "github.com/ory/hydra-client-go/v2"
	"github.com/stretchr/testify/assert"
)

func Test_service_MakeChallenge(t *testing.T) {
	redirectTo := "some other url"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blob, err := json.Marshal(client.NewOAuth2RedirectTo(redirectTo))
		assert.NoError(t, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(blob)
		assert.NoError(t, err)
	}))
	impl, err := New(Config{
		HydraAdminUrl: server.URL,
	})
	assert.NoError(t, err)
	redirectAddr, err := impl.MakeChallenge(context.Background(), "", "challenge")
	assert.NoError(t, err)
	assert.Equal(t, redirectTo, redirectAddr)
	server.Close()
}
