package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LoggerMidddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request - Method:%s | Status: %d | Duration: %v", c.Request.Method,c.Writer.Status(),duration)
	}
}

func main(){
      //create a new Gin router
	router := gin.Default()

	//use our custom logger middleware
	router.Use(LoggerMidddleware())

	//Define a route for the root url

	router.GET("/",func (c *gin.Context)  {
		c.String(200, "Hello, world!")
	})

	//run the server
	router.Run(":8080")

}