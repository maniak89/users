package handlers

import (
	_ "embed"
	"net/http"
)

//go:embed login.html
var loginPage []byte

func LoginGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(loginPage)
}
