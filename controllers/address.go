package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shash-786/EcommerceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("Empty ProductID given")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("No ProductID Available"))
		}

		user_obj_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("Cannot Form ObjectID from Hex")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		empty_addresses := make([]models.Address, 0)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		filter := bson.D{{"_id", user_obj_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{"address", empty_addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			log.Println("DeleteAddress mongoCollection.UpdateOne Error")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		c.IndentedJSON(http.StatusOK, "Successfully Deleted Address")
	}
}
