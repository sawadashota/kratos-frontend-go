package salary

import (
	"net/http"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/sawadashota/kratos-frontend-go/middleware"
	"github.com/sawadashota/kratos-frontend-go/x"
	"github.com/sirupsen/logrus"
)

var (
	salaryHTML *x.HTMLTemplate
)

func init() {
	compileTemplate()
}

func compileTemplate() {
	box := x.NewBox(packr.New("salary", "./templates"))
	salaryHTML = box.MustParseHTML("salary", "layout.html", "salary.html")
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
	sub.Use(h.r.Middleware().Authorize)

	sub.HandleFunc("/my/salary", h.RenderSalary).Methods(http.MethodGet)
}

func (h *Handler) RenderSalary(w http.ResponseWriter, r *http.Request) {
	type salary struct {
		Annual   uint
		Currency string
	}
	htmlValues := struct {
		LogoutURL string
		Salary    salary
	}{
		LogoutURL: h.c.KratosLogoutURL(),
		Salary: salary{
			Annual:   5000000000000000,
			Currency: "JPY",
		},
	}

	if err := salaryHTML.Render(w, &htmlValues); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
