package middleware

import (
	"net/http"

	"github.com/sawadashota/kratos-gin-frontend/x/jwttoken"

	"github.com/gin-gonic/gin"
	"github.com/sawadashota/kratos-gin-frontend/driver"
)

type Middleware struct {
	d driver.Driver
}

func New(d driver.Driver) *Middleware {
	return &Middleware{
		d: d,
	}
}

func (m *Middleware) JWTProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwttoken.ParseRequest(c.Request, m.d.Configuration().JWKsURL())

		if err != nil {
			m.d.Registry().Logger().Infof("fail to authentication: %s", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !token.Valid {
			m.d.Registry().Logger().Info("token is invalid")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := jwttoken.ParseTokenClaims(token.Claims)
		if err != nil {
			m.d.Registry().Logger().Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		m.d.Registry().Logger().Debugf("claims raw: %v", token.Claims)
		m.d.Registry().Logger().Debugf("claims: %v", claims)

		c.Next()
	}
}
