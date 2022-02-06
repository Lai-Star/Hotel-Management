package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/middleware"
)

func main() {

	fmt.Println("restaurant management backend")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.Use(middleware.Authendication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)

}
