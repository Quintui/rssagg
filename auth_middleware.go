package main

import (
	"net/http"

	"github.com/quintui/rssagg/internal/auth"
	"github.com/quintui/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg apiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}

		handler(w, r, user)
	}
}
