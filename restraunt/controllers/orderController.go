package controllers

import (
	"context"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "orders") 
var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

func GetOrders() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		result , err := orderCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch orders"})
			return
		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode orders"})
			return
		}
		c.JSON(200, allOrders)
		defer cancel()
		
	}
}

func GetOrder() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		foodId := c.Param("food_id")
		var order models.Food

		if foodId == "" {
			c.JSON(400, gin.H{"error": "food_id is required"})
			return
		}

		err := foodCollection.FindOne(ctx, bson.M{"order_id": foodId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch order"})
			return
		}
		c.JSON(200, order)
		defer cancel()


		
	}
}

func CreateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}