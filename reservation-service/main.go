package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"reservation-service/handlers"
	"reservation-service/repositories"
	"time"

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

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// NoSQL: Initialize Reservation Repository store
	store, err := repositories.New(storeLogger, timeoutContext)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.CloseSession(timeoutContext)
	store.CreateTables()

	store.Ping()

	//Initialize the handler and inject said logger
	reservationHandler := handlers.NewReservationsHandler(logger, store)

	router := gin.Default()
	//router.Use(accommodationHandler.MiddlewareContentTypeSet(), accommodationHandler.CORSMiddleware())
	router.Use(reservationHandler.CORSMiddleware())

	rg1 := router.Group("/", reservationHandler.CORSMiddleware())
	rg2 := router.Group("/host", reservationHandler.CORSMiddleware(), reservationHandler.Authorize("HOST"))
	rg3 := router.Group("/guest", reservationHandler.CORSMiddleware(), reservationHandler.Authorize("GUEST"))

	rg1.GET("/accommodation/:location/:accomm_id", reservationHandler.GetResByAccommodation)
	rg1.GET("/accommodation-price/:accomm_id", reservationHandler.GetPriceByAccommodation)
	rg1.GET("/accommodation/price-variation/:location/:accomm_id", reservationHandler.GetPriceVariationByAccommId)
	rg1.POST("/check-price/", reservationHandler.CheckPrice)
	rg1.GET("/availability/:location/:accomm_id", reservationHandler.GetAvailability)
	rg1.POST("/search/", reservationHandler.SearchAccommodations)
	rg1.POST("/check-pending/", reservationHandler.GetPendingReservationsByUser)

	rg2.POST("/accommodation-price/", reservationHandler.PostPrice)
	rg2.POST("/accommodation/price-variation/", reservationHandler.CreatePriceVariation)
	rg2.POST("/accommodation/price-variation/delete", reservationHandler.DeletePriceVariation)
	rg2.POST("/availability/create/:id", reservationHandler.CreateAvailability)
	rg2.POST("/availability/delete", reservationHandler.DeleteAvailability)

	rg3.POST("/", reservationHandler.PostReservation)
	rg3.POST("/delete/", reservationHandler.DeleteReservation)

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
