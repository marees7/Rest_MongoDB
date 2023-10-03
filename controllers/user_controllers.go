package controllers

import (
	"api/auth"
	"api/models"
	"api/responses"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type Server struct {
	Client *mongo.Client
	Router *gin.Engine
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var user models.User
	_ = json.NewDecoder(ctx.Request.Body).Decode(&user)
	collection := s.Client.Database("sample").Collection("users")
	c, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(c, user)
	json.NewEncoder(ctx.Writer).Encode(result)
}

func (s *Server) UserLogin(ctx *gin.Context) {
	var login models.Login
	var user models.User
	_ = json.NewDecoder(ctx.Request.Body).Decode(&login)
	collection := s.Client.Database("sample").Collection("users")
	c, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(c, models.User{Email: login.Email, Password: login.Password}).Decode(&user)
	if err != nil {
		json.NewEncoder(ctx.Writer).Encode(models.LoginResponse{})
		return
	}
	loginResponse, err := auth.CreateToken(user.UserID, user.Email)

	if err != nil {
		json.NewEncoder(ctx.Writer).Encode(models.LoginResponse{})
		return
	}
	json.NewEncoder(ctx.Writer).Encode(loginResponse)
}

func (s *Server) GetLoggedInUser(ctx *gin.Context) {
	authEmail, err := auth.ExtractTokenID(ctx.Request)
	if err != nil {
		responses.Error(ctx.Writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if err != nil {
		responses.Error(ctx.Writer, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	user := models.User{}
	userGotten, err := user.WhoAmI(s.Client, authEmail)
	if err != nil {
		responses.Error(ctx.Writer, http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(200, userGotten)
}
