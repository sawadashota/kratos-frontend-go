package account

import (
	"net/http"

	"github.com/sawadashota/kratos-gin-frontend/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/gobuffalo/packr/v2"
	"github.com/sawadashota/kratos-gin-frontend/x/htmlpackr"
)

var (
	homeHTML *htmlpackr.HTMLTemplate
)

func init() {
	compileTemplate()
}

func compileTemplate() {
	box := htmlpackr.New(packr.New("account", "./templates"))
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

func (h *Handler) RenderHome(w http.ResponseWriter, _ *http.Request) {
	htmlValues := struct {
		LogoutURL string
	}{
		LogoutURL: h.c.KratosLogoutURL(),
	}

	if err := homeHTML.Render(w, &htmlValues); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
