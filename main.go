package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/controllers"
	"github.com/shash-786/EcommerceBackend/database"
	"github.com/shash-786/EcommerceBackend/middleware"
	"github.com/shash-786/EcommerceBackend/routes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	var client *mongo.Client = database.Client
	var prodCollection *mongo.Collection = database.ProductData(client, "Product")
	var userCollection *mongo.Collection = database.UserData(client, "User")

	app := controllers.NewApplication(prodCollection, userCollection)

	router.GET("/user/addtocart", app.AddToCart())
	router.GET("/user/removefromcart", app.RemoveItemFromCart())
	router.GET("/user/getitemsfromcart", app.GetItemFromCart())
	router.GET("/user/buycart", app.BuyFromCart())
	router.GET("/user/instantbuy", app.InstantBuy())

	collections, err := client.Database("Ecommerce").ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		log.Println("Error listing collections:", err)
		return
	}
	fmt.Println("Collections in the database:")
	for _, collection := range collections {
		fmt.Println(collection)
	}

	log.Fatal(router.Run(":" + port))
}
