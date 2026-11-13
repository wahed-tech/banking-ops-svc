package api

import (
	"net/http"

	"github.com/wahed-tech/banking-ops-svc/auth"
	"github.com/wahed-tech/go-lib/v2/server"
)

// NewRouter returns the GraphQL surface consumed by the ops back-office portal,
// mounted on top of go-lib's health router so /healthz keeps working.
func NewRouter(pingHandler *PingHandler, authenticator *auth.Authenticator, graphQL http.Handler, allowedOrigins []string) http.Handler {
	mux := server.NewHealthRouter()
	mux.HandleFunc("GET /v1/ping", pingHandler.Ping())
	mux.Handle("GET /graphql", authenticator.Require(graphQL))
	mux.Handle("POST /graphql", authenticator.Require(graphQL))

	return server.WithTracing(withCORS(allowedOrigins, mux))
}
