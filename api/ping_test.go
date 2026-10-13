//go:build unit

package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(_ context.Context) error { return f.err }

func TestPingHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		wantBody   PingResponse
	}{
		{
			name:       "postgres reachable",
			wantStatus: http.StatusOK,
			wantBody:   PingResponse{Status: "ok", Postgres: "ok"},
		},
		{
			name:       "postgres unreachable",
			pingErr:    errors.New("connection refused"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   PingResponse{Status: "degraded", Postgres: "unreachable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/ping", http.NoBody)
			NewRouter(NewPingHandler(fakePinger{err: tt.pingErr})).ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			var got PingResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
			assert.Equal(t, tt.wantBody, got)
		})
	}
}

func TestRouterKeepsHealthz(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", http.NoBody)
	NewRouter(NewPingHandler(fakePinger{})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}
