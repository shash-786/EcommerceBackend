package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/controllers"
)

func UserRoutes(incoming_routes *gin.Engine) {
	incoming_routes.POST("/user/signup", controllers.Signup())
	incoming_routes.POST("/user/login", controllers.Login())
	incoming_routes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incoming_routes.GET("/user/productview", controllers.SearchProduct())
	incoming_routes.GET("/user/search", controllers.SearchProductByQuery())
	incoming_routes.POST("/user/add_address", controllers.AddAddress())
	incoming_routes.PUT("/user/edit_home_address", controllers.EditHomeAddress())
	incoming_routes.PUT("/user/edit_work_address", controllers.EditWorkAddress())
	incoming_routes.POST("/user/delete_address", controllers.DeleteAddress())
}
