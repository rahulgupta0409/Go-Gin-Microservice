package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/user/configs"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	UserType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// func GenerateAllTokens(email string, firstname string, lastname string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
// 	claims := &SignedDetails{
// 		Email:     email,
// 		FirstName: firstname,
// 		LastName:  lastname,
// 		Uid:       uid,
// 		UserType:  userType,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(5)).Unix(),
// 		},
// 	}
// 	token, err := jwt.NewWithClaims(jwt.SigningMethodH256, claims).SignedString([]byte(SECRET_KEY))
// 	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodH256, claims).SignedString([]byte(SECRET_KEY))
// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}
// 	return token, refreshToken, err
// }

// func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&SignedDetails{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			return []byte(SECRET_KEY), nil
// 		},
// 	)
// 	if err != nil {
// 		msg = err.Error()
// 		return
// 	}
// 	claims, ok := token.Claims.(*SignedDetails)
// 	if !ok {
// 		msg = fmt.Sprintf("the token is invalid")
// 		msg = err.Error()
// 	}
// 	if claims.ExpiresAt < time.Now().Unix() {
// 		msg = fmt.Sprintf("token is expired")
// 		msg = err.Error()
// 		return
// 	}
// 	return claims, msg
// }

// func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var updateObj primitive.D
// 	updateObj = append(updateObj, bson.E{"token", signedToken})
// 	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshToken})

// 	ModifiedDate, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 	updateObj = append(updateObj, bson.E{"modifieddate", ModifiedDate})

// 	upsert := true
// 	filter := bson.M{"userid": userId}
// 	opt := options.UpdateOptions{
// 		Upsert: &upsert,
// 	}

// 	_, err := usersCollection.UpdateOne(
// 		ctx,
// 		filter,
// 		bson.D{
// 			{"$set", updateObj}},
// 		&opt,
// 	)

// 	defer cancel()

// 	if err != nil {
// 		log.Panic(err)
// 		return
// 	}
// 	return
// }

func GenerateAllTokens(email string, firstName string, lastName string, userType string, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserId:    userId,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refreshtoken", signedRefreshToken})

	ModifiedDate, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"modifieddate", ModifiedDate})

	upsert := true
	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{"$set", updateObj}},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}
