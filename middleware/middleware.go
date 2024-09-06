package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/tokens"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			log.Println("No Token Found")
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("No Token Found"))
		}

		claims, msg := tokens.ValidateToken(ClientToken)
		if msg != "" {
			log.Println(msg)
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("error Validating Token"))
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
