package controllers

import (
	"context"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")


func GetMenus() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the menus"})
			return
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
		defer cancel()
		
	}
}

func GetMenu() gin.HandlerFunc{
	// Write your code here
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id":menuId}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while fetching the menu"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
} 

func CreateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"error occurred while binding the menu"})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil{
			msg := "menu item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"error occurred while binding the menu"})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}
		
		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil{

			if !inTimeSpan(*menu.Start_Date, *menu.End_Date,time.Now()){
				c.JSON(http.StatusBadRequest, gin.H{"error":"start date must be before end date"})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

			if menu.Name != nil{
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Description != nil{
				updateObj = append(updateObj, bson.E{Key: "description", Value: menu.Description})
			}
			if menu.Category != nil{
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{Key: "$set", Value: updateObj},
				},
				&opt,

			)

			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while updating the menu"})
				defer cancel()
				return
			}
			defer cancel()
			if result.MatchedCount == 1{
				c.JSON(http.StatusOK, gin.H{"message":"menu item was updated","menu": result})
			}else{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"menu item was not updated"})
			}

		}
	}
}

func inTimeSpan(time1, time2, check time.Time) bool {
	return check.After(time1) && check.Before(time2)
}