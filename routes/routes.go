package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/controllers"
)

func UserRoutes(incoming_routes *gin.Engine) {
	incoming_routes.POST("/user/signup", controllers.SignUp())
	incoming_routes.POST("/user/login", controllers.Login())
	incoming_routes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incoming_routes.GET("/user/productview", controllers.SearchProduct())
	incoming_routes.GET("/user/search", controllers.SearchProductByQuery())
}
