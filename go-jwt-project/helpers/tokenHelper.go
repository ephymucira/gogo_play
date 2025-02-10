package helpers

import(
	"context"
	"fmt"
	"log"
	"os"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ephymucira/go-jwt-project/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignedDetails struct {
	Email                  string
	Uid                    string
	First_name             string
	Last_name              string
	User_type              string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, uid string, userType string) (signedToken string,signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		Uid: uid,
		First_name: firstName,
		Last_name: lastName,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},

	}
	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
	}
	refreshToken,err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
	}

	return token,refreshToken,err
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if claims, ok := token.Claims.(*SignedDetails); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func UpdateAllTokens(signedToken string, signedRefreshToken string ,userId string){

    var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)

	var uptdateObj primitive.D

	uptdateObj = append(uptdateObj, bson.E{"token", signedToken})
	uptdateObj = append(uptdateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at ,_:= time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	uptdateObj = append(uptdateObj, bson.E{"updated_at", Updated_at})

	upsert := true
	filter := bson.M{"uid": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx, 
		filter, 
		bson.D{{"$set", uptdateObj}}, 
		&opt,
	)

	defer cancel()

	if err!= nil{
		log.Panic(err)
		return
	}
	
}