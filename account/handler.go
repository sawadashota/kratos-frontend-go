package account

import (
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/sawadashota/kratos-frontend-go/middleware"
	"github.com/sawadashota/kratos-frontend-go/x"
	"github.com/sirupsen/logrus"
)

var (
	homeHTML *x.HTMLTemplate
)

func init() {
	compileTemplate()
}

func compileTemplate() {
	box := x.NewBox(packr.New("account", "./templates"))
	homeHTML = box.MustParseHTML("home", "layout.html", "home.html")
}

type Handler struct {
	r Registry
	c Configuration
}

type Registry interface {
	Logger() logrus.FieldLogger
	Middleware() *middleware.Middleware
}

type Configuration interface {
	KratosLogoutURL() string
	JWKsURL() string
}

func New(r Registry, c Configuration) *Handler {
	return &Handler{
		r: r,
		c: c,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	sub := router.NewRoute().Subrouter()
	sub.Use(h.r.Middleware().JWTProtection)

	sub.HandleFunc("/", h.RenderHome).Methods(http.MethodGet)
}

func (h *Handler) RenderHome(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetClaimsFromContext(r)
	if err != nil {
		h.r.Logger().Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claimsJSON, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		h.r.Logger().Errorf("fail to marshal claims: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	htmlValues := struct {
		LogoutURL  string
		ClaimsJSON string
	}{
		LogoutURL:  h.c.KratosLogoutURL(),
		ClaimsJSON: string(claimsJSON),
	}

	if err := homeHTML.Render(w, &htmlValues); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
