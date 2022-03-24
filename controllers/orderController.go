package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

//get all the order func
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	}
}

//get the single order based on ID
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//create a new order func
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//update order based on ID
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
