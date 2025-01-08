package main

import (
	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	//basic route

	router.GET("/", func(c *gin.Context)  {
		c.String(200, "Hello, World!")
	})

	//route with url parameters

	router.GET("users/:id", func(c *gin.Context)  {
		id := c.Param("id")
		c.String(200, "User ID:"+id)
	})

	//route to query parameters
	router.GET("search/", func(ctx *gin.Context) {
		query := ctx.DefaultQuery("q","default-value")
		ctx.String(200, "Search query:"+query)
	})

	router.Run(":8080")
}