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

// user sends a request for a password Recovery code email
func CreatePasswordRecoveryCode() gin.HandlerFunc {
	return func(c *gin.Context) {

		var foundUser models.User

		var email models.EmailModel

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.ShouldBindJSON(&email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		existingEmail := email.Email
		/* for every new password code request, delete the possibly existing Recovery
		   code for the given email even if the recovery code didn't expire yet */
		res, err := emailVerifCollection.DeleteOne(ctx, bson.M{"verifUsername": existingEmail})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		fmt.Print(res)
		fmt.Print(res.DeletedCount)

		// search User collection for the account with the given email
		var userError = userCollection.FindOne(ctx, bson.M{"email": email.Email}).Decode(&foundUser)
		if userError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Account with this email doesn't exist!"})
			return
		}

		if !foundUser.Is_verified {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email isn't verified yet!"})
			return
		}

		/* create a user object that holds the generated Recovery code, Code creation time
		   and Email of the user that wants to change his password */
		var code = utils.GenerateRandomString(8)
		var userVerif models.UserVerifModel
		userVerif.VerifUsername = email.Email
		userVerif.Code = &code
		t := time.Now().UTC()
		userVerif.Created_at = &t

		// send the Recovery code to the given email
		helper.SendVerifPasswordCode(*userVerif.VerifUsername, *userVerif.Code)

		/* save the 'forgot password' user object into a collection that holds the email Verification
		and password Recovery codes */
		emailVerifCollection.InsertOne(ctx, userVerif)

		c.JSON(http.StatusOK, "")
	}
}

// this handler is called after the user submits his new password with the Recovery code from his email
func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userToChangePass models.ForgotPasswordModel

		if err := c.ShouldBindJSON(&userToChangePass); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// method that checks if the user submitted the wrong code or the code itself expired (after 1 minute)
		checkError := CheckPasswordRecoveryCode(userToChangePass.Email, userToChangePass.Code, ctx)

		if checkError != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": &checkError})
			return
		}

		password := HashPassword(*userToChangePass.Password)

		// if the Recovery code was correct and didn't expire then change the users old password into the new one
		filter := bson.D{{Key: "email", Value: userToChangePass.Email}}
		update := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "password", Value: password},
			},
		}}
		_, errr := userCollection.UpdateOne(ctx, filter, update)
		if errr != nil {
			fmt.Print(errr)
		}

		c.JSON(http.StatusOK, "")

	}
}

func CheckPasswordRecoveryCode(email *string, code *string, ctx context.Context) string {

	// search the user Recovery codes database to see if the given email and code match
	var foundVerifUser models.UserVerifModel
	err := emailVerifCollection.FindOne(ctx, bson.M{"verifUsername": email, "code": code}).Decode(&foundVerifUser)

	if err != nil {
		return "Wrong code!"
	}

	/* for optimal testing purposes the Recovery code expires after just 60 seconds
	   but the Code won't actually be deleted from the database after expiration if,
	   for example, the user requests a Recovery code email but doesn't use it because
	   he changed his mind or suddenly remembered his old password */

	/* TODO see if there is a way to automatically delete these password Recovery codes from
	   the database in case that users request Codes without sending them for checks */
	timeOfCodeCheck := time.Now().UTC()
	timeOfCodeExpiration := foundVerifUser.Created_at.Add(time.Second * 60)

	if timeOfCodeCheck.After(timeOfCodeExpiration) {
		res, errr := emailVerifCollection.DeleteOne(ctx, bson.D{{Key: "verifUsername", Value: &email}, {Key: "code", Value: code}})
		fmt.Print(res)
		fmt.Print(res.DeletedCount)
		if err != nil {
			fmt.Print(errr.Error())
		}
		return "Code expired!"
	}

	return ""
}
