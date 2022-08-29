package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
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

//get all the invoices
func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {

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
