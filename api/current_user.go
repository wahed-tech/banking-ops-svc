package api

import (
	"encoding/json"
	"net/http"

	"github.com/wahed-tech/banking-ops-svc/auth"
)

type currentUserResponse struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func GetCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := auth.CallerFrom(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(currentUserResponse{UserID: caller.UserID, Email: caller.Email})
	}
}
