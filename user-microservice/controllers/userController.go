package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/user/helpers"
	"example.com/user/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"example.com/user/configs"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for email"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for phone"})
			c.Abort()
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		}
		user.CreatedDate = time.Now()
		user.ModifiedDate = primitive.NilObjectID.Timestamp().Local()
		user.ID = primitive.NewObjectID()
		userUUID := uuid.New()
		user.UserId = (userUUID).String()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		if count == 0 {
			resultInsertationNumber, insertErr := userCollection.InsertOne(ctx, user)

			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
				return
			}
			defer cancel()
			c.JSON(http.StatusOK, resultInsertationNumber)
		}
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.LoginDto
		var foundUser models.User
		var userResponse models.UpdateUser

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if foundUser.IsActive == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is deleted"})
		} else {
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect"})
				return
			}
			passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
			defer cancel()
			if passwordIsValid != true {
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			}
			if foundUser.Email == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "please enter an email"})
			}
			token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)
			helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)
			err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&userResponse)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, gin.H{"userid": userResponse.UserId, "firstname": userResponse.FirstName, "lastname": userResponse.LastName,
				"email": userResponse.Email, "phone": userResponse.Phone, "token": userResponse.Token, "refreshtoken": userResponse.RefreshToken,
				"usertype": userResponse.UserType, "isactive": userResponse.IsActive})

		}
	}
}
