package middleware

import (
	"net/http"

	"github.com/sawadashota/kratos-gin-frontend/x/jwttoken"
)

func (m *Middleware) JWTProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwttoken.ParseRequest(r, m.c.JWKsURL())

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

		claims, err := jwttoken.ParseTokenClaims(token.Claims)
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		m.r.Logger().Debugf("claims raw: %v", token.Claims)
		m.r.Logger().Debugf("claims: %v", claims)

		next.ServeHTTP(w, r)
	})
}
