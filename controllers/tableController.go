package controllers

import (
	"context"
	"go-hotel/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.Opencollection(database.Client, "table")

//get all the tables
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := tableCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing tables items"})
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTables)

	}
}

//get the single table by Id
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//create new table api func
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//update table data using ID
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
