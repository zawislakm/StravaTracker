package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

type service struct {
	client *mongo.Client
}

type databaseVariables struct {
	URI    string
	DbName string
}
type indexSetupData struct {
	collectionName string
	indexKeys      bson.D
}

var (
	defaultTimeout = 2 * time.Minute

	dbClient    *service
	dbVariables *databaseVariables

	athletesCollection       = "athletes"
	activitiesCollection     = "activities"
	athleteDataSumCollection = "athleteDataSum"

	indexes = []indexSetupData{
		{
			collectionName: athletesCollection,
			indexKeys: bson.D{
				{Key: "firstname", Value: 1},
				{Key: "lastname", Value: 1},
			},
		},
		{
			collectionName: athleteDataSumCollection,
			indexKeys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "year", Value: 1},
			},
		},
	}
	onceAfterConnection sync.Once
)

func init() {
	dbClient = &service{
		client: nil,
	}
	if dbVariables == nil {
		dbVariables = &databaseVariables{
			URI:    os.Getenv("DB_URI"),
			DbName: os.Getenv("DB_NAME"),
		}
	}
	log.Println("Database variables initialized")
}

func GetDbClient() Service {
	log.Println("Getting database client service")

	if dbClient == nil {
		dbClient = &service{
			client: nil,
		}
	}
	onceAfterConnection.Do(func() {
		if err := dbClient.setupIndexes(); err != nil {
			log.Fatalf("Error setting up indexes: %v", err)
		}
		dbClient.recalculateAllAthletesDataSumUpForYear(time.Now().Format("2006"))
	})
	return dbClient
}

func (s *service) setupIndexes() error {
	log.Println("Setting up indexes")
	for _, index := range indexes {
		log.Println(fmt.Sprintf("Setting up index for collection: %s, with keys: %+v", index.collectionName, index.indexKeys))
		collection, err := s.getCollection(index.collectionName)
		if err != nil {
			return err
		}
		indexModel := mongo.IndexModel{
			Keys:    index.indexKeys,
			Options: options.Index().SetUnique(true),
		}

		_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) getClientConnection(uri string) error {
	if s.isConnected() {
		return nil
	}
	log.Println("Connecting to database: ", uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.
		Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI).SetMaxPoolSize(20).
		SetMinPoolSize(2).
		SetMaxConnIdleTime(defaultTimeout)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	s.client = client
	if !s.isConnected() {
		return fmt.Errorf("failed connecting to database")
	}

	return nil
}

func (s *service) isConnected() bool {
	if s.client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.client.Ping(ctx, nil)
	return err == nil
}

func (s *service) getCollection(collection string) (*mongo.Collection, error) {
	// getCollection gets a MongoDB collection and updates the last activity time
	log.Println("Getting collection: ", collection)
	if err := s.getClientConnection(dbVariables.URI); err != nil {
		return nil, err
	}
	return s.client.Database(dbVariables.DbName).Collection(collection), nil
}

func (s *service) clear() error {
	// function that drops all records from the database
	log.Println("Clearing database")
	return s.client.Database(dbVariables.DbName).Drop(context.Background())
}

func (s *service) close() error {
	// Close manually closes the MongoDB connection
	if !s.isConnected() {
		return nil
	}
	log.Println("Closing MongoDB connection")
	return s.client.Disconnect(context.Background())
}
