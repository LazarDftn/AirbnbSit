package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"profile-service/handlers"
	"profile-service/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	// Initialize context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[accommodation-store] ", log.LstdFlags)

	store, err := repositories.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)

	// NoSQL: Checking if the connection was established
	store.Ping()

	//Initialize the handler and inject said logger
	profileHandler := handlers.NewProfileHandler(logger, store)

	router := gin.Default()
	router.Use(profileHandler.CORSMiddleware())
	router.GET("/", profileHandler.GetAllProfiles)
	router.GET("/:email", profileHandler.GetProfile)
	router.POST("/create", profileHandler.Signup)

	// adding the authorization middleware for Guest and Host, depending on the User action route
	//router.Use(accommodationHandler.Authorize("HOST")).POST("/create", accommodationHandler.PostAccommodation)

	router.Run(":" + port)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

}
