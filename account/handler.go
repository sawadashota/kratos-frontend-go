package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/sawadashota/kratos-gin-frontend/driver"
	"github.com/sawadashota/kratos-gin-frontend/middleware"
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
	router *gin.Engine
	d      driver.Driver
	mw     *middleware.Middleware
}

func New(d driver.Driver, router *gin.Engine, mw *middleware.Middleware) *Handler {
	return &Handler{
		router: router,
		d:      d,
		mw:     mw,
	}
}

func (h *Handler) RegisterRoutes() {
	group := h.router.Group("/", h.mw.JWTProtection())

	group.GET("/", h.RenderHome)
}

func (h *Handler) RenderHome(c *gin.Context) {
	htmlValues := struct {
		LogoutURL string
	}{
		LogoutURL: h.d.Configuration().KratosLogoutURL(),
	}

	if err := homeHTML.Render(c.Writer, &htmlValues); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Status(200)
}
