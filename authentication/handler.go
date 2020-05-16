package authentication

import (
	"net/http"
	"path/filepath"

	"github.com/ory/kratos-client-go/models"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/ory/kratos-client-go/client/common"
	"github.com/sawadashota/kratos-gin-frontend/driver"
	"github.com/sawadashota/kratos-gin-frontend/x/htmlpackr"
)

var (
	signInHTML *htmlpackr.HTMLTemplate
	signUpHTML *htmlpackr.HTMLTemplate
)

func init() {
	compileTemplate()
}

func compileTemplate() {
	box := htmlpackr.New(packr.New("authentication", "./templates"))
	signInHTML = box.MustParseHTML("signin", "layout.html", "signin.html")
	signUpHTML = box.MustParseHTML("signup", "layout.html", "signup.html")
}

type Handler struct {
	router *gin.Engine
	d      driver.Driver
}

func New(d driver.Driver, router *gin.Engine) *Handler {
	return &Handler{
		router: router,
		d:      d,
	}
}

func (h *Handler) RegisterRoutes() {
	group := h.router.Group("auth")

	group.GET("signin", h.RenderSignInForm).Use(h.ValidateRequest)
	group.GET("signup", h.RenderSignUpForm).Use(h.ValidateRequest)
}

func (h *Handler) ValidateRequest(c *gin.Context) {
	requestCode := c.Request.URL.Query().Get("request")
	if requestCode == "" {
		base := filepath.Base(c.Request.URL.Path)
		if base == "login" {
			c.Redirect(http.StatusFound, h.d.Configuration().KratosLoginURL())
		}
		c.Redirect(http.StatusFound, h.d.Configuration().KratosRegistrationURL())
	}
}

func (h *Handler) RenderSignInForm(c *gin.Context) {
	requestCode := c.Request.URL.Query().Get("request")
	params := common.NewGetSelfServiceBrowserLoginRequestParams().WithRequest(requestCode)
	res, err := h.d.Registry().KratosClient().Common.GetSelfServiceBrowserLoginRequest(params)
	if err != nil {
		h.d.Registry().Logger().Errorf("fail to get login request from kratos: %s", err)
		c.Redirect(http.StatusFound, h.d.Configuration().KratosLoginURL())
		return
	}
	if res.Error() == "" {
		h.d.Registry().Logger().Errorf("fail to get login request from kratos: %s", res.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	form := res.GetPayload().Methods["password"].Config

	htmlValues := struct {
		Form *models.LoginRequestMethodConfig
	}{
		Form: form,
	}
	if err := signInHTML.Render(c.Writer, &htmlValues); err != nil {
		h.d.Registry().Logger().Errorf("fail to render HTML: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}

func (h *Handler) RenderSignUpForm(c *gin.Context) {
	requestCode := c.Request.URL.Query().Get("request")
	params := common.NewGetSelfServiceBrowserRegistrationRequestParams().WithRequest(requestCode)
	res, err := h.d.Registry().KratosClient().Common.GetSelfServiceBrowserRegistrationRequest(params)
	if err != nil {
		h.d.Registry().Logger().Errorf("fail to get registration request from kratos: %s", err)
		c.Redirect(http.StatusFound, h.d.Configuration().KratosRegistrationURL())
		return
	}
	if res.Error() == "" {
		h.d.Registry().Logger().Errorf("fail to get registration request from kratos: %s", res.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	form := res.GetPayload().Methods["password"].Config

	htmlValues := struct {
		Form *models.RegistrationRequestMethodConfig
	}{
		Form: form,
	}
	if err := signUpHTML.Render(c.Writer, &htmlValues); err != nil {
		h.d.Registry().Logger().Errorf("fail to render HTML: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(200)
}
