package repositories

import (
	"accommodations-service/domain"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type AccommodationRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*AccommodationRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	return &AccommodationRepo{
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (pr *AccommodationRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (ar *AccommodationRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := ar.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		ar.logger.Println(err)
	}

	// Print available databases
	databases, err := ar.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
	}
	fmt.Println(databases)
}

func (ar *AccommodationRepo) GetAll() (domain.Accommodations, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accommCollection := ar.getCollection()

	var accommodations domain.Accommodations
	accommodationsCursor, err := accommCollection.Find(ctx, bson.M{})
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	if err = accommodationsCursor.All(ctx, &accommodations); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodations, nil
}

func (ar *AccommodationRepo) Insert(accommodation *domain.Accommodation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	accommCollection := ar.getCollection()

	result, err := accommCollection.InsertOne(ctx, &accommodation)
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	ar.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (ar *AccommodationRepo) getCollection() *mongo.Collection {
	accommDatabase := ar.cli.Database("mongoDemo")
	accommCollection := accommDatabase.Collection("accommodations")
	return accommCollection
}
