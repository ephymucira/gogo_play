package controllers

import (
	"context"
	"fmt"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		endIndex := startIndex + recordPerPage
		var total int64
		count, err := foodCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while counting the documents"})
			return
		}
		total = count
		var foods []bson.M
		cursor, err := foodCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the foods"})
			return
		}
		if err = cursor.All(ctx, &foods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while parsing the foods"})
			return
		}
		if endIndex > len(foods) {
			endIndex = len(foods)
		}
		foods = foods[startIndex:endIndex]
		defer cancel()
		c.JSON(http.StatusOK, gin.H{
			"total": total,
			"page": page,
			"recordPerPage": recordPerPage,
			"data": foods,
		})
		defer cancel()
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
	return int(num+math.Copysign(0.5,num))
	

}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output

}

func UpdateFood() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var food models.Food
		var menu models.Menu

		foodId := c.Param("food_id")
		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"error occurred while binding the food"})
			return
		}
		validationErr := validate.Struct(food)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}
		var updateObj primitive.D
		if food.Name != nil{
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}
		if food.Description != nil{
			updateObj = append(updateObj, bson.E{Key: "description", Value: food.Description})
		}
		if food.Price != nil{
			num := toFixed(*food.Price, 2)
			food.Price = &num
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
		}
		if food.Food_image != nil{
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.Food_image})
		}
		if food.Menu_id != nil{
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.Menu_id})
			err := menuCollection.FindOne(ctx, bson.M{"menu_id":*food.Menu_id}).Decode(&menu)
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the menu"})
				return
			}
			if menu.Start_Date != nil && menu.End_Date != nil{
				if !inTimeSpan(*menu.Start_Date, *menu.End_Date,time.Now()){
					c.JSON(http.StatusBadRequest, gin.H{"error":"start date must be before end date"})
					defer cancel()
					return
				}
			}
		}
		food.UpdatedAt ,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})
		upsert := true
		filter := bson.M{"food_id":foodId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while updating the food"})
			defer cancel()
			return
		}
		defer cancel()
		if result.MatchedCount == 1{
			c.JSON(http.StatusOK, gin.H{"message":"food item was updated successfully","food": food})
		}else{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"food item was not updated"})
		}
		

	}
}