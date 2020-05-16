package middleware

import (
	"net/http"
	"path/filepath"
)

func (m *Middleware) ValidateFormRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCode := r.URL.Query().Get("request")
		if requestCode == "" {
			base := filepath.Base(r.URL.Path)
			if base == "signin" {
				http.Redirect(w, r, m.c.KratosLoginURL(), http.StatusFound)
				return
			}
			http.Redirect(w, r, m.c.KratosRegistrationURL(), http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
