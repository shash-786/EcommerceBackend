package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(productCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_query_id := c.Query("id")
		if user_query_id == "" {
			log.Println("User ID is nil")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No User ID Given",
			})
			return
		}

		prod_query_id := c.Query("prod_id")
		if prod_query_id == "" {
			log.Println("Product ID is nil")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No Product ID Given",
			})
			return
		}

		product_obj_id, err := primitive.ObjectIDFromHex(prod_query_id)
		if err != nil {
			log.Println("Cannot Form Object ID From Product Query ID")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.productCollection, app.productCollection, product_obj_id, user_query_id)
		if err != nil {
			log.Println("panic: controllers/cart Database AddtoCart error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.IndentedJSON(http.StatusOK, "Successfully Added Product To Cart")
	}
}

func (app *Application) RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_query_id := c.Query("id")
		if user_query_id == "" {
			log.Println("The User Query ID Is nil")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Empty UserID Given"))
		}

		product_query_id := c.Query("prod_id")
		if product_query_id == "" {
			log.Println("The Product Query ID Is nil")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Empty Product Given"))
		}

		product_obj_id, err := primitive.ObjectIDFromHex(product_query_id)
		if err != nil {
			log.Println("Cannot Form Object ID From Product Query ID")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveProductFromCart(ctx, app.productCollection, app.userCollection, user_query_id, product_obj_id)
		if err != nil {
			log.Println("panic: controllers/cart Database RemoveProductFromCart error")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		c.IndentedJSON(http.StatusOK, "Successfully Removed Item From Cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_query_id := c.Query("prod_id")
		if product_query_id == "" {
			log.Println("No Product Query ID Found")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Empty Product Query ID"))
		}

		product_obj_id, err := primitive.ObjectIDFromHex(product_query_id)
		if err != nil {
			log.Println("Couldn't Convert QueryID to valid ObjID")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		user_query_id := c.Query("id")
		if user_query_id == "" {
			log.Println("The User Query ID Is nil")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Empty UserID Given"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err = database.InstantBuy(ctx, app.productCollection, app.userCollection, user_query_id, product_obj_id)
		if err != nil {
			log.Println("panic: controllers/cart Database InstantBuy error")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "Order (InstaBuy) Successfully Placed")
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_query_id := c.Query("id")
		if user_query_id == "" {
			log.Println("The User Query ID Is nil")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Empty UserID Given"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.productCollection, app.userCollection, user_query_id)
		if err != nil {
			log.Println("Error Processing the payment")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, "Order Placed Successfully")
	}
}
