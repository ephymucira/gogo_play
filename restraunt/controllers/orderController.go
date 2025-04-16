package controllers

import (
	"context"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		var order models.Order
		var table models.Table
		if err := c.BindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(order)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		if order.TableId != "" {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to fetch table"})
				return
			}
		}
		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()
		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			msg := "order item was not created"
			c.JSON(500, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(200, result)
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
		var table models.Table
		var order models.Order

		var updateObj primitive.D
		orderId := c.Param("order_id")
		if orderId == "" {
			c.JSON(400, gin.H{"error": "order_id is required"})
			return
		}
		if err := c.BindJSON(&order); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		if order.TableId != "" {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to fetch table"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableId})
			order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})
			
			upsert := true
			filter := bson.M{"order_id": orderId}
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result,err := orderCollection.UpdateOne(
				ctx,
				filter,
				bson.D{{Key: "$set", Value: updateObj}},
				&opt,
			)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to update order"})
				return
			}
			if result.MatchedCount == 0 {
				c.JSON(404, gin.H{"error": "Order not found"})
				return
			}
			if result.ModifiedCount == 0 {
				c.JSON(404, gin.H{"error": "Order not modified"})
				return
			}
			c.JSON(200, gin.H{"message": "Order updated successfully", "order": result})
			defer cancel()

		}
	}
}

func OrderItemOrderCreator(order models.Order) string{
	order.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()
	return order.OrderId
}