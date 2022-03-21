package controllers

import (
	"context"
	"go-hotel/models"
	"github.com/gin-gonic/gin"
	"time"
)

// var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
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

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

}
