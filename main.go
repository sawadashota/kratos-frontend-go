package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sawadashota/kratos-gin-frontend/account"
	"github.com/sawadashota/kratos-gin-frontend/authentication"
	"github.com/sawadashota/kratos-gin-frontend/driver"
	"github.com/sawadashota/kratos-gin-frontend/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	d := driver.NewDefaultDriver()
	mw := middleware.New(d)
	router := gin.Default()

	authentication.New(d, router).RegisterRoutes()
	account.New(d, router, mw).RegisterRoutes()

	if err := router.Run(fmt.Sprintf(":%d", d.Configuration().Port())); err != nil {
		logrus.Fatalf("fail to staring server: %s", err)
	}
}
