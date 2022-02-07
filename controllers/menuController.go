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

var menuCollection *mongo.Collection = database.Opencollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurd while fetching menu data"})
		}

		var menu []bson.M
		if err := result.All(ctx, &menu); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, menu)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
