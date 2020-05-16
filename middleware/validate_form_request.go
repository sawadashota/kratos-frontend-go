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
			switch base {
			case "settings":
				http.Redirect(w, r, m.c.KratosSettingsURL(), http.StatusFound)
				return
			case "signin":
				http.Redirect(w, r, m.c.KratosLoginURL(), http.StatusFound)
				return
			default:
				http.Redirect(w, r, m.c.KratosRegistrationURL(), http.StatusFound)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
