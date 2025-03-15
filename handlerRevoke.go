package main

import (
	"context"
	"net/http"

	"github.com/chaseplamoureux/gohttpserver/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, 400, "no token provided", err)
	}

	err = cfg.db.RevokeRefreshToken(context.Background(), bearerToken)
	if err != nil {
		errorResponse(w, 500, "could not revoke token", err)
	}
	responseJSON(w, 204, nil)
}