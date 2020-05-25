package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/sawadashota/kratos-frontend-go/internal/jwt"
)

func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		token, err := extractJWTToken(r)
		if err != nil {
			m.r.Logger().Info(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		type Input struct {
			Method string   `json:"method"`
			User   string   `json:"user"`
			Path   []string `json:"path"`
			Token  string   `json:"token"`
		}
		input := struct {
			Input Input `json:"input"`
		}{
			Input: Input{
				Method: r.Method,
				Path:   path,
				Token:  token,
			},
		}

		b, err := json.Marshal(input)
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest(http.MethodPost, m.c.OPAPolicyURL(), bytes.NewBuffer(b))
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cl := http.DefaultClient
		res, err := cl.Do(req)
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if res.StatusCode >= 300 {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type respBody struct {
			DecisionID string `json:"decision_id"`
			Result     struct {
				Allow bool `json:"allow"`
				Token struct {
					Payload jwt.Claims `json:"payload"`
				} `json:"token"`
			} `json:"result"`
		}
		var rb respBody
		if err := json.NewDecoder(res.Body).Decode(&rb); err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !rb.Result.Allow {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if rb.Result.Token.Payload.IsExpired() {
			http.Redirect(w, r, m.c.KratosLogoutURL(), http.StatusFound)
			return
		}

		m.r.Logger().Debug(rb)
		r = SetClaimsToContext(r, &rb.Result.Token.Payload)

		next.ServeHTTP(w, r)
	})
}

func extractJWTToken(r *http.Request) (string, error) {
	s := strings.Split(r.Header.Get("authorization"), " ")
	if len(s) != 2 {
		return "", errors.New("authorization header is not found")
	}
	return s[1], nil
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
