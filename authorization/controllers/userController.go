package controllers

import (
	"auth/database"
	helper "auth/helpers"
	"auth/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var emailVerifCollection *mongo.Collection = database.OpenCollection(database.Client, "emailVerif")
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
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.Is_verified = false

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		usernameCount, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
		}

		if emailCount > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		if usernameCount > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exists"})
			return
		}

		user.ID = primitive.NewObjectID()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.ID.Hex())
		user.Token = &token
		user.Refresh_token = &refreshToken

		var code = utils.GenerateRandomString(8)
		var userVerif models.UserVerifModel
		userVerif.VerifUsername = user.Username
		userVerif.Code = &code

		emailVerifCollection.InsertOne(ctx, userVerif)
		helper.SendVerifEmail(user, code)

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.ID.Hex())
		helper.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.ID.Hex()}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func VerifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.UserVerifModel
		var foundUser models.UserVerifModel

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		err := emailVerifCollection.FindOne(ctx, bson.M{"verifUsername": user.VerifUsername, "code": user.Code}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This user is already verified or isn't registered yet!"})
			return
		}

		filter := bson.D{{"username", user.VerifUsername}}
		update := bson.D{{"$set",
			bson.D{
				{"is_verified", true},
			},
		}}
		_, errr := userCollection.UpdateOne(ctx, filter, update)
		if errr != nil {
			fmt.Print(errr)
		}

		emailVerifCollection.DeleteOne(ctx, bson.D{{"verifUsername", user.VerifUsername}})

		c.JSON(http.StatusOK, "")
	}
}

func CreatePasswordRecoveryCode() gin.HandlerFunc {
	return func(c *gin.Context) {

		var foundUser models.User
		var email string

		if err := c.ShouldBindJSON(email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user = userCollection.FindOne(ctx, bson.M{"Email": foundUser.Email}).Decode(&foundUser)
		defer cancel()
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User with this email address doesn't exist!"})
			return
		}

		/*for every new password code request, delete the existing code even
		if it didn't expire yet*/
		emailVerifCollection.DeleteOne(ctx, bson.D{{"verifUsername", email}})

		var code = utils.GenerateRandomString(8)
		var userVerif models.UserVerifModel
		userVerif.VerifUsername = &email
		userVerif.Code = &code
		t := time.Now().UTC()
		userVerif.Created_at = &t

		emailVerifCollection.InsertOne(ctx, userVerif)

		c.JSON(http.StatusOK, "")
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userToChangePass models.ChangePasswordModel

		if err := c.ShouldBindJSON(userToChangePass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		checkError := CheckPasswordRecoveryCode(userToChangePass.Email, userToChangePass.Code)

		if checkError != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": &checkError})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		password := HashPassword(*userToChangePass.Password)

		filter := bson.D{{"email", userToChangePass.Email}}
		update := bson.D{{"$set",
			bson.D{
				{"password", password},
			},
		}}
		_, errr := userCollection.UpdateOne(ctx, filter, update)
		if errr != nil {
			fmt.Print(errr)
		}

		c.JSON(http.StatusOK, "")

	}
}

func CheckPasswordRecoveryCode(email *string, code *string) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var foundVerifUser models.UserVerifModel
	err := emailVerifCollection.FindOne(ctx, bson.M{"verifUsername": email, "code": code}).Decode(&foundVerifUser)

	if err != nil {
		return "Wrong code!"
	}

	t1 := time.Now().UTC()
	t2 := foundVerifUser.Created_at.Add(time.Second * 60)

	if t1.After(t2) {
		emailVerifCollection.DeleteOne(ctx, bson.D{{"verifUsername", email}, {"code", code}})
		return "Code expired!"
	}

	return ""
}
