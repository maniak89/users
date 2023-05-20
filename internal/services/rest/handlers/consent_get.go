package handlers

import (
	"errors"
	"net/http"

	serviceImpl "users/internal/services/service"
)

func ConsentGet(impl *serviceImpl.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		redirectTo, err := impl.Consent(ctx, r.URL.Query().Get("consent_challenge"))
		if err != nil {
			if errors.Is(err, serviceImpl.ErrNoUser) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
	}
}
