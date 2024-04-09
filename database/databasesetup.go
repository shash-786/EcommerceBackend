package controller

import "github.com/gin-gonic/gin"

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(entered_password string, password_in_db string) (bool, string) {
	return
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
