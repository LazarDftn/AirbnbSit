package repositories

import (
	"context"
	"fmt"
	"log"
	"profile-service/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ProfileRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger) (*ProfileRepo, error) {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://accommodations_db:27017/"))
	if err != nil {
		return nil, err
	}

	return &ProfileRepo{
		cli:    client,
		logger: logger,
	}, nil
}

func DBinstance() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://auth_db:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

// Disconnect from database
func (pr *ProfileRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *ProfileRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (pr *ProfileRepo) GetAll() (domain.Users, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	var profiles []domain.User
	profilesCursor, err := profileCollection.Find(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	if err = profilesCursor.All(ctx, &profiles); err != nil {
		pr.logger.Println(err)
		return nil, err
	}
	return profiles, nil
}

func (pr *ProfileRepo) GetProfile(email string) (*domain.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	var profile domain.User

	err := profileCollection.FindOne(ctx, bson.M{"email": email}).Decode(&profile)
	if err != nil {
		pr.logger.Println(err)
		return &profile, err
	}
	return &profile, nil
}

func (pr *ProfileRepo) Insert(user domain.User) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	res, err := profileCollection.InsertOne(ctx, user)

	return res, err

}

func (pr *ProfileRepo) CheckUsernameExists(username string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	count, err := profileCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		pr.logger.Println(err)
		return "error"
	}

	if count > 0 {
		return "username exists"
	}

	return ""
}

func (pr *ProfileRepo) Delete(id primitive.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	profileCollection := pr.getCollection()

	_, err := profileCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})

	return err
}

func (ar *ProfileRepo) getCollection() *mongo.Collection {
	accommDatabase := ar.cli.Database("mongoDemo")
	accommCollection := accommDatabase.Collection("profile")
	return accommCollection
}
