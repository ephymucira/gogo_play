package routes

import (
	"github.com/gin-gonic/gin"
	"golang-restraunt-management/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())
}