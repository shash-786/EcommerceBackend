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
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("User Id Field Empty")
			c.AbortWithStatus(http.StatusBadRequest)
		}

		user_obj_id, _ := primitive.ObjectIDFromHex(user_id)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var new_user_address models.Address
		new_user_address.Address_ID = primitive.NewObjectID()

		if err := c.BindJSON(&new_user_address); err != nil {
			log.Println("Couldn't bind the address properly")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		filter_match := bson.D{
			{
				Key: "$match", Value: bson.D{
					primitive.E{Key: "_id", Value: user_obj_id},
				},
			},
		}

		unwind := bson.D{
			{
				Key: "$unwind", Value: bson.D{
					primitive.E{Key: "$path", Value: "$address"},
				},
			},
		}

		group := bson.D{
			{
				Key: "$group", Value: bson.D{
					primitive.E{Key: "_id", Value: user_obj_id},
					primitive.E{Key: "total", Value: bson.D{
						primitive.E{Key: "$sum", Value: 1},
					}},
				},
			},
		}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, group})
		defer cursor.Close(ctx)

		if err != nil {
			log.Println("Aggregate Failed")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		var aggregate_query_results []bson.M
		if err = cursor.All(ctx, &aggregate_query_results); err != nil {
			log.Println("Couldn't Decode the cursor")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		var size int32
		for _, primitive_map := range aggregate_query_results {
			count := primitive_map["total"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{{Key: "_id", Value: user_obj_id}}

			update := bson.D{
				{
					Key: "$push", Value: bson.D{
						primitive.E{Key: "address", Value: new_user_address},
					},
				},
			}

			if _, err = UserCollection.UpdateOne(ctx, filter, update); err != nil {
				log.Println("Error Updating the User")
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			c.IndentedJSON(http.StatusOK, "Successfully Added An Address")

		} else {
			log.Println("Already Max Addresses Present")
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("Empty ProductID given")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("No ProductID Available"))
		}

		user_obj_id, _ := primitive.ObjectIDFromHex(user_id)

		var edit_address models.Address
		if err := c.BindJSON(&edit_address); err != nil {
			log.Println("Couldn't bind the address properly")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: user_obj_id}}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					primitive.E{Key: "$address.0.house", Value: edit_address.House},
					primitive.E{Key: "$address.0.street", Value: edit_address.Street},
					primitive.E{Key: "$address.0.city", Value: edit_address.City},
					primitive.E{Key: "$address.0.pincode", Value: edit_address.Pincode},
				},
			},
		}

		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			log.Println("EditAddress UpdateOne Error")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully Edited Home Address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("Empty ProductID given")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("No ProductID Available"))
		}

		user_obj_id, _ := primitive.ObjectIDFromHex(user_id)

		var edit_address models.Address
		if err := c.BindJSON(&edit_address); err != nil {
			log.Println("Couldn't bind the address properly")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		filter := bson.D{{Key: "_id", Value: user_obj_id}}
		update := bson.D{
			{
				Key: "$set", Value: bson.D{
					primitive.E{Key: "$address.1.house", Value: edit_address.House},
					primitive.E{Key: "$address.1.street", Value: edit_address.Street},
					primitive.E{Key: "$address.1.city", Value: edit_address.City},
					primitive.E{Key: "$address.1.pincode", Value: edit_address.Pincode},
				},
			},
		}

		if _, err := UserCollection.UpdateOne(ctx, filter, update); err != nil {
			log.Println("EditAddress UpdateOne Error")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.IndentedJSON(http.StatusOK, "Successfully Edited Office Address")
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
