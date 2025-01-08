package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMidddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method:%s | Status: %d | Duration: %v", c.Request.Method, c.Writer.Status(), duration)
	}
}

//custom auth middleware

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apikey := c.GetHeader("X-API-KEY")
		if apikey != "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthrized"})
		}
		c.Next()
	}
}

func main() {
	//create a new Gin router
	router := gin.Default()

	//use our custom logger middleware
	// router.Use(LoggerMidddleware())

	//Define a route for the root url

	// router.GET("/",func (c *gin.Context)  {
	// 	c.String(200, "Hello, world!")
	// })

	authGroup := router.Group("/api")
	authGroup.Use(AuthMiddleware())
	{
		authGroup.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Authenticated and authorized!"})
		})
	}

	//run the server
	router.Run(":8080")

}
