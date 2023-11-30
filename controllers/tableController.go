package controllers

import (
	"context"
	"fmt"
	"go-hotel/database"
	"go-hotel/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.Opencollection(database.Client, "table")

//get all the tables
func GetTables(c *gin.Context) {
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

//get the single table by Id
func GetTable(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	tableId := c.Param("table_id")
	var table models.Table
	err := orderCollection.FindOne(ctx, bson.M{"order_id": tableId}).Decode(&table)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the tables"})
	}
	c.JSON(http.StatusOK, table)

}

//create new table api func
func CreateTable(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var table models.Table
	if err := c.BindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(table)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	table.ID = primitive.NewObjectID()
	table.Table_id = table.ID.Hex()

	result, inserterr := tableCollection.InsertOne(ctx, table)
	if inserterr != nil {
		msg := fmt.Sprintf("Table item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result)

}

//update table data using ID
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		tableid := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
		}
		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		filter := bson.M{"table_id": tableid}

		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)

		if err != nil {
			msg := fmt.Sprintf("table item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}
