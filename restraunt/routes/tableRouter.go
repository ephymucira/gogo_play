package routes

import (
	"github.com/gin-gonic/gin"
	"golang-restraunt-management/controllers"
)

func invoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/tables", controllers.GetTables())
	incomingRoutes.GET("/tables/:table_id", controllers.GetTable())
	incomingRoutes.POST("/tables", controllers.CreateTable())
	incomingRoutes.PATCH("/tables/:table_id", controllers.UpdateTable())
}