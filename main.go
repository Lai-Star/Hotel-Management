package main

import (
	"fmt"
	"os"

	"go-hotel/database"
	"go-hotel/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.Opencollection(database.Client, "food")

func main() {

	fmt.Println("restaurant management backend")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)

}
