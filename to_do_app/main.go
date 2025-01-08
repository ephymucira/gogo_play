package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct{
	gorm.Model
	Title        string `json:"tittle"`
	Description  string `json:"description"`
}

func main(){
	router := gin.Default()

	//connect to the database
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	//Auto-migrate the Todo model to create the table

	db.AutoMigrate(&Todo{})

    //Route to create a todo
	router.POST("/todos", func(ctx *gin.Context) {
		var todo Todo
		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(400, gin.H{"error":"Invalid JSON data"})
			return
		}
		db.Create(&todo)

		ctx.JSON(200,todo)
	})
	//Route to get all todos

	router.GET("/todos", func(ctx *gin.Context) {
		var todos []Todo
		
		//retrieve all todos from the database
		db.Find(&todos)

		ctx.JSON(200, todos)
	})

	//route to get a specific todo

	router.PUT("/todos/:id", func(ctx *gin.Context) {
		var todo Todo
		todoID := ctx.Param("id")
	
		// Retrieve the Todo from the database
		result := db.First(&todo, todoID)
		if result.Error != nil {
			ctx.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
	
		// Bind JSON input to updatedTodo
		var updatedTodo Todo
		if err := ctx.ShouldBindJSON(&updatedTodo); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid JSON input", "details": err.Error()})
			return
		}
	
		// Update the Todo fields
		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
	
		// Save the updated Todo to the database
		if err := db.Save(&todo).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to update Todo"})
			return
		}
	
		ctx.JSON(200, gin.H{"message": "Todo updated successfully", "todo": todo})
	})
	

	//Route to delete a Todo by ID
	router.DELETE("/todos/:id", func(ctx *gin.Context) {
		var todo Todo
		todoID := ctx.Param("id")
	
		// Retrieve the Todo from the database
		result := db.First(&todo, todoID)
		if result.Error != nil {
			ctx.JSON(404, gin.H{"error": "Todo not found"})
			return
		}
	
		// Delete the Todo from the database
		if err := db.Delete(&todo).Error; err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to delete Todo"})
			return
		}
	
		// Respond with success message
		ctx.JSON(200, gin.H{
			"message": "Todo deleted successfully",
			"todo_id": todoID,
		})
	})
	

	router.Run(":8080")

	
}