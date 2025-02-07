package routes

import(
	"github.com/gin-gonic/gin"
	controller "github.com/ephymucira/go-jwt-project/controllers"
	"github.com/ephymucira/go-jwt-project/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.GET("users", controller.GetUsers())
	incomingRoutes.GET("users/:user_id", controller.GetUser())
}