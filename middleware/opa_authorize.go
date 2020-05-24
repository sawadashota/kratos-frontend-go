package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func (m *Middleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaimsFromContext(r)
		path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		m.r.Logger().Infof("request from %s", claims.Sub)

		type Input struct {
			Method string   `json:"method"`
			User   string   `json:"user"`
			Path   []string `json:"path"`
		}
		input := struct {
			Input Input `json:"input"`
		}{
			Input: Input{
				Method: r.Method,
				User:   claims.Sub,
				Path:   path,
			},
		}

		b, err := json.Marshal(input)
		if err != nil {
			m.r.Logger().Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest(http.MethodPost, "http://opa:8181/v1/data/frontend", bytes.NewBuffer(b))
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

		next.ServeHTTP(w, r)
	})
}
