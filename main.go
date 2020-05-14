package main

import (
	"fmt"
	"net/http"

	"github.com/sawadashota/kratos-gin-frontend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/sawadashota/kratos-gin-frontend/authentication"
	"github.com/sawadashota/kratos-gin-frontend/driver"
	"github.com/sawadashota/kratos-gin-frontend/x/htmlpackr"
	"github.com/sirupsen/logrus"
)

var (
	indexHTML *htmlpackr.HTMLTemplate
)

func init() {
	compileTemplate()
}

func compileTemplate() {
	box := htmlpackr.New(packr.New("HTML templates", "./templates"))
	indexHTML = box.MustParseHTML("index", "layout.html", "index.html")
}

func main() {
	d := driver.NewDefaultDriver()
	mw := middleware.New(d)
	router := gin.Default()

	authentication.New(d, router).RegisterRoutes()

	protected := router.Group("/", mw.JWTProtection())
	protected.GET("/", func(c *gin.Context) {
		if err := indexHTML.Render(c.Writer, gin.H{}); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.Status(200)
	})

	if err := router.Run(fmt.Sprintf(":%d", d.Configuration().Port())); err != nil {
		logrus.Fatalf("fail to staring server: %s", err)
	}
}
