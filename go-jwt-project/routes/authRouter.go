package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/ephymucira/go-jwt-project/controllers"
	

)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/register", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}