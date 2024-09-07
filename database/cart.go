package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/shash-786/EcommerceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, product_obj_id primitive.ObjectID, usr_id string) error {
	usr_obj_id, err := primitive.ObjectIDFromHex(usr_id)
	if err != nil {
		log.Println("database/cart : Cannot create Object ID")
		return err
	}

	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": product_obj_id})
	if err != nil {
		log.Println("Can't find Product")
		return err
	}

	var product_cart []models.ProductUser
	err = searchfromdb.All(ctx, &product_cart)
	if err != nil {
		log.Println("Can't decode Product")
		return err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: usr_obj_id},
	}

	update := bson.D{
		primitive.E{Key: "$push", Value: bson.D{
			primitive.E{Key: "user_cart", Value: bson.D{
				primitive.E{Key: "$each", Value: product_cart},
			}},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("AddProdcutToCart: Error in Updating the User")
		return err
	}
	return nil
}

func RemoveProductFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, usr_id string, product_obj_id primitive.ObjectID) error {
	usr_obj_id, err := primitive.ObjectIDFromHex(usr_id)
	if err != nil {
		log.Println("database/cart : Cannot create Object ID")
		return err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: usr_obj_id},
	}

	update := bson.D{
		primitive.E{Key: "$pull", Value: bson.D{
			primitive.E{Key: "user_cart", Value: bson.D{
				primitive.E{Key: "_id", Value: product_obj_id},
			}},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error in Removing Prod from cart/ updating user cart")
		return err
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, usr_id string) error {
	user_obj_id, err := primitive.ObjectIDFromHex(usr_id)
	if err != nil {
		log.Println("database/cart : Cannot create Object ID")
		return err
	}

	var user models.User
	var order models.Order

	err = userCollection.FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: user_obj_id},
	}).Decode(&user)

	if len(user.User_Cart) == 0 {
		return errors.New("UserCart Empty")
	}

	if err != nil {
		log.Println("Error in decoding the user")
		return err
	}

	order.Order_ID = primitive.NewObjectID()
	order.Ordered_At = time.Now()

	payment := models.Payment{
		Digital: false,
		COD:     true,
	}

	order.Payment_Method = payment

	filter_match := bson.D{
		primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "_id", Value: user_obj_id},
		}},
	}

	unwind := bson.D{
		primitive.E{Key: "$unwind", Value: bson.D{
			primitive.E{Key: "$path", Value: "$user_cart"},
		}},
	}

	grouping := bson.D{
		primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: user_obj_id},
			primitive.E{Key: "total_price", Value: bson.D{
				primitive.E{Key: "$sum", Value: "$user_cart.price"},
			}},
		}},
	}

	cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
	if err != nil {
		log.Println("Error in the Aggregate Pipeline")
		return err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Println("Error in the Cursor Decoding")
		return err
	}

	var total_price int
	for _, val := range results {
		price := val["total_price"]
		total_price = price.(int)
	}
	order.Price = total_price //ASSIGNMENT OF THE TOTAL PRICE ONTO THE ORDER STRUCT
	order.Order_Cart = make([]models.ProductUser, len(user.User_Cart))
	copy(order.Order_Cart, user.User_Cart)

	filter := bson.D{
		primitive.E{Key: "$_id", Value: user_obj_id},
	}

	update := bson.D{
		primitive.E{Key: "$push", Value: bson.D{
			primitive.E{Key: "$orders", Value: order},
		}},
	}

	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		log.Println(err)
		return err
	}

	new_empty_user_cart := make([]models.ProductUser, 0)
	update1 := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "$user_cart", Value: new_empty_user_cart},
		}},
	}

	if _, err := userCollection.UpdateOne(ctx, filter, update1); err != nil {
		log.Println("Cannot Update the Purchase")
		return err
	}
	return nil
}

func InstantBuy(ctx context.Context, userCollection, prodCollection *mongo.Collection, usr_id string, product_obj_id primitive.ObjectID) error {
	user_obj_id, err := primitive.ObjectIDFromHex(usr_id)
	if err != nil {
		log.Println("database/cart : Cannot create Object ID")
		return err
	}

	var product_to_buy models.ProductUser
	err = prodCollection.FindOne(ctx, bson.D{
		primitive.E{Key: "_id", Value: product_obj_id},
	}).Decode(&product_to_buy)

	if err != nil {
		log.Println("InstantBuy: Product Not Found")
		return err
	}

	var order models.Order
	order.Order_ID = primitive.NewObjectID()
	order.Ordered_At = time.Now()

	payment := models.Payment{
		Digital: false,
		COD:     true,
	}

	order.Payment_Method = payment
	order.Price = product_to_buy.Price
	order.Order_Cart = make([]models.ProductUser, 0)
	order.Order_Cart = append(order.Order_Cart, product_to_buy)

	filter := bson.D{
		primitive.E{Key: "_id", Value: user_obj_id},
	}

	update := bson.D{
		primitive.E{Key: "$push", Value: bson.D{
			primitive.E{Key: "$orders", Value: order},
		}},
	}

	if _, err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		log.Println("Cannot Update the Purchase")
		return err
	}
	return nil
}
