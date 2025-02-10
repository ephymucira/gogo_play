package middleware

import (
	"net/http"
	helper "github.com/ephymucira/go-jwt-project/helpers"
	"github.com/gin-gonic/gin"
)

// Authenticate is a middleware that checks if the user is authenticated
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Next()
		
		
	}
}