package controllers

import (
	"go-hotel/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.Opencollection(database.Client, "table")

//get all the tables
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {

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
