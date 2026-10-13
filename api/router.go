package api

import (
	"net/http"

	"github.com/wahed-tech/go-lib/v2/server"
)

// NewRouter returns the REST surface consumed by the ops back-office portal,
// mounted on top of go-lib's health router so /healthz keeps working.
func NewRouter(pingHandler *PingHandler) http.Handler {
	mux := server.NewHealthRouter()
	mux.HandleFunc("GET /v1/ping", pingHandler.Ping())

	return server.WithTracing(mux)
}
