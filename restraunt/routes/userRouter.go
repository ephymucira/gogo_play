package routes

import(
	"github.com/gin-gonic/gin"
	"golang-restraunt-management/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	invoiceRoutes.POST("/users/signup", controllers.signup())
	incomingRoutes.POST("/users/login", controllers.login())
}