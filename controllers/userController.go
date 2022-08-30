package controllers

import (
	"context"
	"go-hotel/database"
	"go-hotel/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// user collection
var userCollection *mongo.Collection = database.Opencollection(database.Client, "user")

// var validate = validator.New()

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordperpage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordperpage < 1 {
			recordperpage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 10
		}

		startindex := (page - 1) * recordperpage
		startindex, err = strconv.Atoi(c.Query("startIndex"))

		matchstage := bson.D{{"$match", bson.D{{}}}}
		projectstage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_details", bson.D{{"$slice", []interface{}{"$data", startindex, recordperpage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchstage, projectstage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, result)
	}
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
	return func(c *gin.Context) {
		//convert the JSON data coming from postman to something that golang understands

		//validate the data based on user struct

		//you'll check if the email has already been used by another user

		//hash password

		//you'll also check if the phone no. has already been used by another user

		//create some extra details for the user object - created_at, updated_at, ID

		//generate token and refersh token (generate all tokens function from helper)

		//if all ok, then you insert this new user into the user collection

		//return status OK and send the result back
	}
}

//login api
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		//convert the login data from postman which is in JSON to golang readable format

		//find a user with that email and see if that user even exists

		//then you will verify the password

		//if all goes well, then you'll generate tokens

		//update tokens - token and refersh token

		//return statusOK
	}
}

//helper func
func HashPassword(password string) string {
	return password
}

//verify password func
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	return true, providedPassword
}
