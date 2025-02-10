package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	database "github.com/ephymucira/go-jwt-project/database"
	helper "github.com/ephymucira/go-jwt-project/helpers"
	"github.com/ephymucira/go-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.client,"user")
// var userCollection *mongo.Collection = database.OpenCollection(database.DBinstance(),"user")
var validate = validator.New()

// func HashPassword()

func VerifyPassword(userPassword string, providedPassword string, c *gin.Context) (bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Password is incorrect"
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
			cancel()
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
			cancel()
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken , _ := helper.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id, user.User_type)
         
		user.Token = token
		user.Refresh_token  = refreshToken
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
			cancel()
			return
		}

		user.Password = string(hashedPassword)
		user.User_id = strconv.FormatInt(time.Now().Unix(), 10)

		result, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
			cancel()
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
		defer cancel()

	}

}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel()
			return
		}
		err:= userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is incorrect."})
			return
		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password, c)

        defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email does not exist"})
		}
		token, refreshToken,_ := helper.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id, foundUser.User_type)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while logging in"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": foundUser})

	}

	
}

// func GetUsers() gin.HandlerFunc{
// 	return func (c *gin.Context) {
// 		helper.CheckUserType(c, "ADMIN"); err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		recordPerPage,err := strconv.Atoi(c.Query("recordPerPage"))

// 		if err != nil or recordPerPage < 1 {
// 			recordPerPage = 10
// 		}
// 		page, err1 := strconv.Atoi(c.Query("page"))
// 		var users []models.User
// 		cursor, err := userCollection.Find(ctx, bson.M{})
// 		if err != nil or page < 1 {
// 			page = 1

// 		}

// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err  = strconv.Atoi(c.Query("startIndex"))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			defer cancel()
// 			return
// 		}

// 		if err = cursor.All(ctx, &users); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching users"})
// 			defer cancel()
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": users})
// 		defer cancel()
		
// 	}
// }

func GetUsers() gin.HandlerFunc{
	return func (c *gin.Context) {
		helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage,err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		var users []models.User
		cursor, err := userCollection.Find(ctx, bson.M{})
		if err != nil || page < 1 {
			page = 1

		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err  = strconv.Atoi(c.Query("startIndex"))
		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "count", Value: 1},
				{Key: "data", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}
		}

		result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		val allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": allusers[0]})

		
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user models.User

		err:= userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})

	}
}

