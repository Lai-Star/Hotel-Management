package middleware

import "github.com/gin-gonic/gin"

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// clienttoken := c.Request.Header.Get("token")

	}
}
