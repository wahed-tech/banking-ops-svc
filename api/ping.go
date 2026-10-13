package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/wahed-tech/go-lib/v2/log"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type PingHandler struct {
	db pinger
}

type PingResponse struct {
	Status   string `json:"status"`
	Postgres string `json:"postgres"`
}

func NewPingHandler(db pinger) *PingHandler {
	return &PingHandler{db: db}
}

// Ping reports whether the service can reach its Postgres, so a dev deploy can be
// smoke-tested end to end through the ingress.
func (h *PingHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, resp := http.StatusOK, PingResponse{Status: "ok", Postgres: "ok"}
		if err := h.db.Ping(r.Context()); err != nil {
			log.Error("postgres ping failed", "", err)
			status, resp = http.StatusServiceUnavailable, PingResponse{Status: "degraded", Postgres: "unreachable"}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
