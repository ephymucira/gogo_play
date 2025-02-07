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