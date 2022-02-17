package controllers

import "github.com/gin-gonic/gin"

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

}
