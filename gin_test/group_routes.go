package main

import (
	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	public :=  router.Group("/public")
	{
		public.GET("/info", func(ctx *gin.Context) {
			ctx.String(200,"Public information")
		})
		public.GET("/products", func(ctx *gin.Context) {
			ctx.String(200, "Public Product list")
		})
	}

	private := router.Group("/private")
	private.Use(AuthMiddleware())
	{
		private.GET("/data", func(ctx *gin.Context) {
			ctx.String(200,"Private data accessible after authentication")
		})
		private.POST("/create",func(ctx *gin.Context) {
			ctx.String(200, "Create a new resource")
		})
	}

	router.Run(":8080")
}