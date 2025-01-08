package main

import (
	"github.com/gin-gonic/gin"

)

//user related controller
type UserController struct{}

// GetUserInfo is a controller method to get user information

func (uc *UserController) GetUserInfo(c *gin.Context){
	userID := c.Param("id")

	//fetch data from the database but for now we shall return a simpe json response
	c.JSON(200, gin.H{"id":userID,"name":"Ephy Mucira","email":"ephymucira@gmail.com"})
}

func main(){
	router := gin.Default()

	userController := &UserController{}

	//Route using the controller
	router.GET("/users/:id", userController.GetUserInfo)

	router.Run(":8080")
}
