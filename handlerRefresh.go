package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chaseplamoureux/gohttpserver/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, 400, "no token provided", err)
	}
	type response struct {
		Token string `json:"token"`
	}
	tokenInfo, err := cfg.db.GetRefreshToken(context.Background(), bearerToken)
	if err != nil {
		errorResponse(w, 401, "token invalid. Unauthorized", err)
	}
	fmt.Println(tokenInfo.ExpiresAt)
	fmt.Println()
	if tokenInfo.ExpiresAt.Before(time.Now().UTC()) {
		errorResponse(w, 401, "token expired. Unauthorized", err)
	}
	if tokenInfo.RevokedAt.Valid && tokenInfo.RevokedAt.Time.Before(time.Now()) {
		errorResponse(w, 401, "token expired. Unauthorized", err)
	}
	newAccessToken, err := auth.MakeJWT(tokenInfo.UserID, cfg.jwt,time.Hour)
	if err != nil {
		errorResponse(w, 500, "Could not create new access token", err)
	}
	responseJSON(w, 200, response{Token: newAccessToken})
}
