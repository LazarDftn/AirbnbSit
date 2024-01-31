package controllers

import (
	"auth/database"
	helper "auth/helpers"
	"auth/models"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/d-vignesh/go-jwt-auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var emailVerifCollection *mongo.Collection = database.OpenCollection(database.Client, "emailVerif")
var validate = validator.New()
var profileAddress string
var resAddress string
var accommAddress string

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func SetAddress() {

	envFile, err := godotenv.Read(".env")

	if err != nil {
		fmt.Println(err.Error())
	}

	profileAddress = envFile["PROFILE_ADDRESS"]
	resAddress = envFile["RESERVATION_ADDRESS"]
	accommAddress = envFile["ACCOMMODATION_ADDRESS"]
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

		unhashedPassword := *user.Password
		password := HashPassword(*user.Password)
		user.Password = &password

		file, err := os.Open("blacklist.txt")

		// Error finding the blacklist.txt file no matter what path we try to use

		if err != nil {
			log.Panic(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var fileLines []string

		for scanner.Scan() {
			fileLines = append(fileLines, scanner.Text())
		}

		file.Close()

		var found string

		for _, line := range fileLines {
			if strings.Contains(unhashedPassword, line) {
				found = line
			}
		}

		if found != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This password is on Blacklist, please change it"})
			return
		}

		usernameCount, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
			return
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
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Username, *user.User_type)
		user.Token = &token
		user.Refresh_token = &refreshToken

		/*
			user.Is_verified = false
			var code = utils.GenerateRandomString(8)
			var userVerif models.UserVerifModel
			userVerif.VerifUsername = user.Username
			userVerif.Code = &code

			emailVerifCollection.InsertOne(ctx, userVerif)
			helper.SendVerifEmail(user, code)
		*/

		var profile models.Profile

		profile.Username = user.Username
		profile.Email = user.Email
		profile.First_name = user.First_name
		profile.Last_name = user.Last_name
		profile.Address = user.Address
		profile.User_type = user.User_type
		profile.ID = primitive.NewObjectID()

		jsonProfile, errr := json.Marshal(&profile)

		if errr != nil {
			fmt.Println(errr)
		}

		SetAddress()

		fmt.Println(profileAddress)

		requestBody := bytes.NewReader(jsonProfile)

		req, err := http.NewRequestWithContext(ctx, "POST", profileAddress+"create", requestBody)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		client := http.Client{Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		}}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register account right now!"})
			return
		} else {
			if resp.StatusCode == 418 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "This username already exists!"})
				return
			}
			if resp.StatusCode != 200 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register account right now!"})
				return
			}
		}

		// Read the response body
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register account right now!"})
			return
		}

		defer resp.Body.Close()

		var credentials models.UserCredentialsModel

		credentials.ID = profile.ID
		credentials.Email = user.Email
		credentials.Password = user.Password
		credentials.Is_verified = true // Change to false later
		var empty = ""
		credentials.Token = &empty
		credentials.Refresh_token = &empty

		_, insertErr := userCollection.InsertOne(ctx, credentials)
		if insertErr != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "")
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var foundUser models.UserCredentialsModel
		var userLogin models.LoginModel
		var foundProfile models.Profile

		if err := c.BindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email is incorrect"})
			return
		}

		if !foundUser.Is_verified {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "To sign in you need to verify your account!"})
			return
		}

		passwordIsValid, _ := VerifyPassword(*userLogin.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password is incorrect"})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		SetAddress()

		req, err := http.NewRequest(http.MethodGet, profileAddress+foundUser.ID.Hex(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving profile"})
			return
		}

		res, errClient := http.DefaultClient.Do(req)

		if errClient != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving profile"})
			return
		}
		if res.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving profile"})
			return
		}

		err = json.NewDecoder(res.Body).Decode(&foundProfile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving profile"})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundProfile.Username, *foundProfile.User_type)
		helper.UpdateAllTokens(token, refreshToken, foundProfile.ID.Hex())
		//err = userCollection.FindOne(ctx, bson.M{"user_id": foundProfile.ID.Hex()}).Decode(&foundProfile)
		foundProfile.Token = &token
		foundProfile.Refresh_token = &refreshToken

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundProfile)
	}
}

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var profile models.Profile

		if err := c.BindJSON(&profile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(profile)
		fmt.Println(profile.Username)
		fmt.Println(&profile.ID)
		fmt.Println(profile.Email)

		jsonProfile, err := json.Marshal(profile)

		fmt.Println(jsonProfile)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		requestBody := bytes.NewReader(jsonProfile)

		req, err := http.NewRequest(http.MethodPost, resAddress+"check-pending/", requestBody)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting profile"})
			return
		}

		res, errClient := http.DefaultClient.Do(req)

		if errClient != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error!"})
			return
		}
		if res.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You have pending reservations!"})
			return
		}

		req, err = http.NewRequest(http.MethodDelete, accommAddress+profile.ID.Hex(), nil)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting profile"})
			return
		}

		res, errClient = http.DefaultClient.Do(req)

		if errClient != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error!"})
			return
		}
		if res.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting your accommodations!"})
			return
		}

		req, err = http.NewRequest(http.MethodDelete, profileAddress+profile.ID.Hex(), nil)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting profile"})
			return
		}

		res, errClient = http.DefaultClient.Do(req)

		if errClient != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error!"})
			return
		}
		if res.StatusCode != 200 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while deleting your accommodations!"})
			return
		}

		result, errr := userCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: profile.ID}})

		if errr != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while deleting your profile!"})
			return
		}

		if result.DeletedCount > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}

		/*fmt.Println(result)

		if errr != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while deleting your profile!"})
			return
		}

		c.JSON(http.StatusOK, "account deleted successfully!")*/
	}
}

func EditAccount(c *gin.Context) {

	id := c.Param("id")

	var user models.UserEdit
	var foundUser models.UserCredentialsModel

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userCollection.FindOne(c, bson.M{"_id": objectId}).Decode(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving profile"})
		return
	}

	if user.NewPassword != "" || user.Email != "" {

		passwordIsValid, _ := VerifyPassword(user.NewPassword, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Old password is incorrect!"})
			return
		}

		if user.Email != "" {

			foundEmail := userCollection.FindOne(c, bson.M{"email": user.Email})

			if foundEmail != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "User with this email already exists!"})
				return
			}
		}
	}

	var userToEdit models.Profile
	userToEdit.Address = &user.Address
	userToEdit.First_name = &user.First_name
	userToEdit.Last_name = &user.Last_name
	userToEdit.Username = &user.Username
	userToEdit.Email = &user.Email

	jsonProfile, errr := json.Marshal(userToEdit)

	if errr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errr.Error()})
		return
	}

	requestBody := bytes.NewReader(jsonProfile)

	req, err := http.NewRequest(http.MethodPut, profileAddress+id, requestBody)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting profile"})
		return
	}

	res, errClient := http.DefaultClient.Do(req)

	if errClient != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error!"})
		return
	}
	if res.StatusCode == 418 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Your new username already exists!"})
		return
	}
	if res.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error while editing your profile!"})
		return
	}

	if user.NewPassword != "" {

		newPassword := HashPassword(user.NewPassword)

		user.NewPassword = newPassword

		userCollection.UpdateOne(c, bson.M{"_id": objectId}, bson.M{
			"$set": bson.M{
				"password": user.NewPassword,
			}})

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return
		}
	}

	if user.Email != "" {
		userCollection.UpdateOne(c, bson.M{"_id": objectId}, bson.M{
			"$set": bson.M{
				"email": user.Email,
			}})

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, "")
			return
		}
	}

	c.JSON(http.StatusOK, "Profile edited successfully!")
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
