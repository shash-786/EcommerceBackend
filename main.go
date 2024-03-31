package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/controllers"
	"github.com/shash-786/EcommerceBackend/database"
	"github.com/shash-786/EcommerceBackend/middleware"
	"github.com/shash-786/EcommerceBackend/models"
	"github.com/shash-786/EcommerceBackend/routes"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.UserRouter(routes)
	router.Use(middleware.Authentication())

	log.Fatal(router.Run())
}
