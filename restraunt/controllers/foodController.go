package controllers

import (
	"context"
	"fmt"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}
  
func GetFood() gin.HandlerFunc{
	return func(c *gin.Context){

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id":foodId}).Decode(&food)

		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the fod "})
		}

		c.JSON(http.StatusOK, food)
		
		
	}
}
var validate = validator.New()

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel  = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"error occurred while binding the food"})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id":food.Menu_id}).Decode(&menu)

		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the menu"})
			return
		}
		food.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.UpdatedAt ,_= time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result ,insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil{
			msg := fmt.Sprintf("food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)


	}
}

func round(num float64) int {
	

}

func toFixed(num float64, precision int) float64 {
	

}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}