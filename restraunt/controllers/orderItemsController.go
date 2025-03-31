package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderItems() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		
	}
}

func GetOrderItem() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func GetOrderItemByOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func ItemsByOrder(id string) (orderItems []primitive.M, err error){
	// Write your code here
	return
}