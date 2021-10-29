package apiserver

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/reqww/go-rest-api/internal/app/auth"
	"net/http"
)

func (s *server) SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestId, id)))
	})
}

func (s *server) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.IsAuthenticated(r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userId, _ := claims["userId"]
			u, err := s.store.User().FindById(int(userId.(float64)))

			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))

		} else {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	})
}
