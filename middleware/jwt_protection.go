package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/sawadashota/kratos-frontend-go/internal/jwt"
)

func (m *Middleware) JWTProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.r.JWTParser().ParseRequest(r)

		if err != nil {
			m.r.Logger().Infof("fail to authentication: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !token.Valid {
			m.r.Logger().Info("token is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := m.r.JWTParser().ParseTokenClaims(token.Claims)
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = SetClaimsToContext(r, claims)

		next.ServeHTTP(w, r)
	})
}

const contextClaimsKey = "claims"

func SetClaimsToContext(r *http.Request, claims *jwt.Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), contextClaimsKey, claims))
}

func GetClaimsFromContext(r *http.Request) (*jwt.Claims, error) {
	claims, ok := r.Context().Value(contextClaimsKey).(*jwt.Claims)
	if !ok {
		return nil, errors.New("request context doesn't have claims")
	}
	return claims, nil
}
