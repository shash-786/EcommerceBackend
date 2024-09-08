package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shash-786/EcommerceBackend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Claims struct {
	Email     string
	Firstname string
	Lastname  string
	Uid       string
	jwt.StandardClaims
}

var (
	secret_key                   = os.Getenv("SECRET")
	SECRET_KEY                   = []byte(secret_key)
	UserData   *mongo.Collection = database.UserData(database.Client, "User")
)

func TokenGenerate(email, firstname, lastname, uid string) (signedtoken, signedrefreshtoken string, err error) {
	claim := &Claims{
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}

	refreshclaim := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(SECRET_KEY)
	if err != nil {
		log.Panicln(err)
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaim).SignedString(SECRET_KEY)
	if err != nil {
		log.Panicln(err)
		return "", "", err
	}

	return token, refreshtoken, err
}

func ValidateToken(signedtoken string) (claims *Claims, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		msg = "Invalid Token"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is Expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedtoken, signedrefreshtoken, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	filter := bson.M{"user_id": userid}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				primitive.E{Key: "token", Value: signedtoken},
				primitive.E{Key: "refresh_token", Value: signedrefreshtoken},
				primitive.E{Key: "updated_at", Value: updated_at},
			},
		},
	}

	_, err := UserData.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Panicln(err)
	}

}
