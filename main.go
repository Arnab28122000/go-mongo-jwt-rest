package main

import (
	"context"
	"log"
	"os"

	"example.com/mongo-apis/controllers"
	"example.com/mongo-apis/database"
	"example.com/mongo-apis/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontroller controllers.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	mongoclient, err := database.DBinstance()
	if err != nil {
		log.Fatal(err)
	}
	usercollection = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(usercollection, ctx)
	usercontroller = controllers.New(userservice)
	server = gin.Default()
}

//v1/user
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basepath)

	log.Fatal(server.Run(":" + port))

}
