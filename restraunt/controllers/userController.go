package controllers

import (

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		
	}
}

func GetUser() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		
	}
}

func signup() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func login() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}


func HashPassword(password string) string{
	// Write your code here
	return ""
}

func VerifyPassword(hashedPassword string, password string) (bool,string){
	// Write your code here
	return false,""
}