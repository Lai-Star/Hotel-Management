package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//get all the order items func
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//get the single order item by ID
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ItemByOrder(id string) (orderItems []primitive.M, err error) {}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//update orderitem based on ID
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//create order item 
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
