package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/mongo-apis/database"
	helper "example.com/mongo-apis/helpers"
	"example.com/mongo-apis/middleware"
	"example.com/mongo-apis/models"
	"example.com/mongo-apis/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserService services.UserService
}

func New(userservice services.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check = false
	}
	return check, msg
}

func (uc *UserController) Signup(ctx *gin.Context) {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		defer cancel()
		return
	}

	count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	count, err = userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		return
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(c, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()
	ctx.JSON(http.StatusOK, resultInsertionNumber)
}

func (uc *UserController) Login(ctx *gin.Context) {
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var foundUser models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
	err = userCollection.FindOne(c, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, foundUser)
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := uc.UserService.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	username := ctx.Param("name")
	// if err := helper.MatchUserTypeToUid(ctx, username); err != nil {
	// 	ctx.JSON(http.StatusBadGateway, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
	user, err := uc.UserService.GetUser(&username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	err := helper.CheckUserType(ctx, "ADMIN")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))

	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(ctx.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(ctx.Query("startIndex"))

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "unable to get start index",
		})
		defer cancel()
		return
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
	result, err := userCollection.Aggregate(c, mongo.Pipeline{
		matchStage, groupStage, projectStage})
	defer cancel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
	}
	var allusers []bson.M
	if err = result.All(ctx, &allusers); err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, allusers[0])

	// users, err := uc.UserService.GetAll()
	// if err != nil {
	// 	ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	// 	return
	// }
	// ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var username string = ctx.Param("name")
	err := uc.UserService.DeleteUser(&username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	authroute := rg.Group("/auth")
	userroute.Use(middleware.Authenticate())
	authroute.POST("/signup", uc.Signup)
	authroute.POST("/login", uc.Login)
	userroute.GET("/get/:name", uc.GetUser)
	userroute.GET("/getall", uc.GetAll)
	userroute.PATCH("/update", uc.UpdateUser)
	userroute.DELETE("/delete/:name", uc.DeleteUser)
}
