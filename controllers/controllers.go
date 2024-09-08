package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shash-786/EcommerceBackend/database"
	"github.com/shash-786/EcommerceBackend/models"
	"github.com/shash-786/EcommerceBackend/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate                            = validator.New()
	UserCollection    *mongo.Collection = database.UserData(database.Client, "User")
	ProductCollection *mongo.Collection = database.ProductData(database.Client, "Product")
)

func HashPassword(password string) string {
	hash_pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hash_pass)
}

func VerifyPassword(entered_password string, password_in_db string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password_in_db), []byte(entered_password))
	valid := true
	msg := ""
	if err != nil {
		valid = false
		msg = "Authentication Failed No Access"
	}

	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		ValidationErr := validate.Struct(user)
		if ValidationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ValidationErr,
			})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User Already Exists",
			})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{
			"phone": user.Phone,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Phone Number Already Exists",
			})
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		// TODO: Implement For Refresh and Tokens
		token, refresh_token, err := tokens.TokenGenerate(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		if err != nil {
			log.Println("Error in token generation")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		user.Token = &token
		user.Refresh_Token = &refresh_token

		user.User_Cart = make([]models.ProductUser, 0)
		user.Address_Detail = make([]models.Address, 0)
		user.Orders = make([]models.Order, 0)

		if _, inserterr := UserCollection.InsertOne(ctx, user); inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "not created",
			})
			return
		}
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var input_user models.User
		var database_user models.User

		if err := c.BindJSON(&input_user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": input_user.Email}).Decode(&database_user)
		if err != nil {
			log.Println("Email or Password Incorrect")
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}

		ValidPass, msg := VerifyPassword(*input_user.Password, *database_user.Password)
		if !ValidPass {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": msg,
			})
			return
		}

		// TODO: UPDATE TOKEN LOGIC

		c.JSON(http.StatusFound, database_user)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			log.Println("Admin product viewer error")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		product.Product_ID = primitive.NewObjectID()

		if _, err := ProductCollection.InsertOne(ctx, product); err != nil {
			log.Println("Admin Insert Product Error!")
			c.JSON(http.StatusBadGateway, gin.H{
				"error": err,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "Product Inserted successfully!",
		})
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// RETURNS ALL THE PRODUCTS PRESENT

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			log.Println("Error in Finding Products")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		products_slice := make([]models.Product, 0)

		err = cursor.All(ctx, &products_slice)
		if err != nil {
			log.Println("Error in copying the Products to product_slice")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		defer cursor.Close(ctx)

		// Redundant Error Check
		if cursor.Err() != nil {
			log.Println("Error Iteration of the Cursor")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		c.IndentedJSON(http.StatusOK, products_slice)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

		product_query_name := c.Query("name")
		if product_query_name == "" {
			log.Println("Not given any Product query name")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("no product name given"))
		}

		searchproductslice := make([]models.Product, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchproductdb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": product_query_name}})
		if err != nil {
			log.Println("Error in Fetching products")
			_ = c.AbortWithError(http.StatusNotFound, err)
		}

		defer searchproductdb.Close(ctx)

		err = searchproductdb.All(ctx, searchproductslice)
		if err != nil {
			log.Println("Error in copying the Products to product_slice")
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		// Redundant Error Check
		if searchproductdb.Err() != nil {
			log.Println("Error Iteration of the Cursor")
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		c.IndentedJSON(http.StatusOK, searchproductslice)
	}
}
