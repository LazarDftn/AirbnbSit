package main

import (
	"log"
	"os"
	"os/signal"
	"reservation-service/handlers"
	"reservation-service/repositories"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	//Initialize the logger we are going to use, with prefix and datetime for every log
	logger := log.New(os.Stdout, "[product-api] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[reservation-store] ", log.LstdFlags)

	// NoSQL: Initialize Reservation Repository store
	store, err := repositories.New(storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseSession()
	store.CreateTables()

	//Initialize the handler and inject said logger
	reservationHandler := handlers.NewReservationsHandler(logger, store)

	router := gin.Default()
	//router.Use(accommodationHandler.MiddlewareContentTypeSet(), accommodationHandler.CORSMiddleware())
	router.POST("/", reservationHandler.PostReservation)
	router.GET("/accommodation/{accomm_id}", reservationHandler.GetResByAccommodation)

	router.POST("/accommodation-price/", reservationHandler.PostPrice)
	router.GET("/accommodation-price/{accomm_id}", reservationHandler.GetPriceByAccommodation)

	// adding the authorization middleware for Guest and Host, depending on the User action route
	// router.Use(accommodationHandler.Authorize("HOST")).POST("/create", accommodationHandler.PostAccommodation)

	router.Run(":" + port)

	/*server := http.Server{
		Addr:         ":" + port,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	logger.Println("Server listening on port", port)
	//Distribute all the connections to goroutines
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()*/

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	logger.Println("Received terminate, graceful shutdown", sig)

	//Try to shutdown gracefully
	//if server.Shutdown(timeoutContext) != nil {
	//	logger.Fatal("Cannot gracefully shutdown...")
	//}

}
