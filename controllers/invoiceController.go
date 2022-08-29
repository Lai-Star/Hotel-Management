package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-hotel/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.Opencollection(database.Client, "invoice")

//get all the invoices
func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoice items"})
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)

	}
}

//get invoice function based on ID
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// create new invoice func
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// update invoice func
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
