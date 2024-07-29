package httpserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/cronnoss/bookshop-home-task/internal/app/common/server"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func (h HTTPServer) CheckAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthorizationHeader)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			server.BadRequest("validate-token", err, w, r)
			return
		}
		if user.Username == "" {
			server.InternalError("invalid-token", nil, w, r)
			return
		}
		if !user.Admin {
			server.Unauthorised("not-admin", nil, w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}

func (h HTTPServer) CheckAuthorizedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthorizationHeader)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			server.BadRequest("validate-token", err, w, r)
			return
		}
		if user.Username == "" {
			server.InternalError("invalid-token", nil, w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}
