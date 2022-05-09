package controllers

import (
	"context"
	"go-hotel/database"
	"go-hotel/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// user collection
var userCollection *mongo.Collection = database.Opencollection(database.Client, "user")

// var validate = validator.New()

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

//get single user
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		c.JSON(http.StatusOK, user)
	}
}

//signup api
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

//login api
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

//helper func
func HashPassword(password string) string {
	return password
}

//verify password func
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	return true, providedPassword
}
