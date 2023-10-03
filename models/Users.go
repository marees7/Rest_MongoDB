package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"github.com/google/uuid"
)

type Server struct {
	Client *mongo.Client
	Router *gin.Engine
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	UserID    uuid.UUID          `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiryDate  int64  `json:"expiry_date"`
}

func (u *User) WhoAmI(db *mongo.Client, authEmail string) (*User, error) {
	collection := db.Database("sample").Collection("users")
	c, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(c, User{Email: authEmail}).Decode(&u)
	if err != nil {
		return &User{}, errors.New("User not found")
	}
	fmt.Println("email ", u.Email)
	return u, nil
}
