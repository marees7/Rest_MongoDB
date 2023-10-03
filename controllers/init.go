package controllers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var s = Server{}

func Initialize() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	s.Client, _ = mongo.Connect(ctx, clientOptions)
	s.Router = gin.Default()
	file, _ := os.Create("log.txt")
	log.SetOutput(file)
	s.InitializeRoutes()
	s.Router.Run(":9191")
}

//
//func Run() {
//}
