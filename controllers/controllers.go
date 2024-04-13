package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shash-786/EcommerceBackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// TODO: CREATE USER AND PRODUCT COLLECTION
var (
	validate          = validator.New()
	UserCollection    *mongo.Client
	ProductCollection *mongo.Client
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

		user.User_Cart = make([]models.ProductUser, 0)
		user.Address_Detail = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

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
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
