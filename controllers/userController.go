package controllers

import (
	"context"
	"fmt"
	"go-hotel/database"
	helper "go-hotel/helpers"
	"go-hotel/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// user collection
var userCollection *mongo.Collection = database.Opencollection(database.Client, "user")

//get the all users
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

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		//convert the JSON data coming from postman to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//validate the data based on user struct
		validationerr := validate.Struct(user)
		if validationerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationerr.Error()})
			return
		}

		//you'll check if the email has already been used by another user
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		//hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		//you'll also check if the phone no. has already been used by another user
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or phone number already exists"})
			return
		}

		//create some extra details for the user object - created_at, updated_at, ID
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//generate token and refersh token (generate all tokens function from helper)
		token, refereshtoken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refereshtoken

		//if all ok, then you insert this new user into the user collection
		resultinsert, insererr := userCollection.InsertOne(ctx, user)
		if insererr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		//return status OK and send the result back
		c.JSON(http.StatusOK, resultinsert)
	}
}

//login api
func Login(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var founduser models.User

	//convert the login data from postman which is in JSON to golang readable format
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//find a user with that email and see if that user even exists
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
		return
	}

	//then you will verify the password
	passwordvalid, msg := VerifyPassword(*user.Password, *founduser.Password)
	defer cancel()

	if passwordvalid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	//if all goes well, then you'll generate tokens
	token, refreshtoken, _ := helper.GenerateAllTokens(*founduser.Email, *founduser.First_name, *founduser.Last_name, founduser.User_id)

	//update tokens - token and refersh token
	helper.UpdateAllTokens(token, refreshtoken, founduser.User_id)

	//return statusOK
	c.JSON(http.StatusOK, founduser)
}

//helper func
func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//verify password func
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}

	return check, msg
}
