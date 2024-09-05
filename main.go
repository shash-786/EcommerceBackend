package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.UserRouter(routes)
	router.Use(middleware.Authentication())

	log.Fatal(router.Run())
}
