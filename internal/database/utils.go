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

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var defaultTimeout = 5 * time.Minute

var dbClient *MongoDBClient
var URI string
var DbName string
var athletesCollection = "athletes"
var activitiesCollection = "activities"
var variableSetUp bool = false

type MongoDBClient struct {
	client       *mongo.Client
	lastActivity time.Time
	timeout      time.Duration
	mu           sync.Mutex
	quit         chan struct{}
}

func GetDbClient() *MongoDBClient {
	log.Println("Getting database client")
	if variableSetUp == false {
		if URI == "" || DbName == "" {
			if err := godotenv.Load(); err != nil {
				log.Fatalf("Error loading .env file: %v", err)
			}
		}

		if URI == "" {
			URI = os.Getenv("DB_URI")
		}
		if DbName == "" {
			DbName = os.Getenv("DB_NAME")
		}
		variableSetUp = true
	}
	if dbClient == nil {
		dbClient = &MongoDBClient{
			client:       nil,
			lastActivity: time.Now(),
			timeout:      defaultTimeout,
			mu:           sync.Mutex{},
			quit:         nil,
		}
	}

	defer func(dbClient *MongoDBClient) {
		err := dbClient.Close()
		if err != nil {
			log.Fatalf("Error closing MongoDB connection: %v", err)
		}
	}(dbClient)
	return dbClient
}

func (service *MongoDBClient) setAthleteIndex() error {
	log.Println("Setting athlete index")
	collection, err := service.getCollection(athletesCollection)
	if err != nil {
		return err
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "firstname", Value: 1},
			{Key: "lastname", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err

	}
	return nil
}

func (service *MongoDBClient) setActivityIndex() error {
	log.Println("Setting activity index")
	collection, err := service.getCollection(activitiesCollection)
	if err != nil {
		return err
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err

	}
	return nil
}

func (service *MongoDBClient) getClientConnection(uri string) error {
	// getClientConnection creates a new MongoDB client with an inactivity timeout
	if service.isConnected() {
		return nil
	}
	log.Println("Connecting to database: ", uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	service.client = client
	if !service.isConnected() {
		return fmt.Errorf("Failed connecting to database")
	}

	service.quit = make(chan struct{})
	//go service.monitorInactivity() removed closing connection due to inactivity
	return nil
}

func (service *MongoDBClient) isConnected() bool {
	if service.client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := service.client.Ping(ctx, nil)
	return err == nil
}

func (service *MongoDBClient) monitorInactivity() {
	// monitorInactivity monitors the client's activity and closes the connection after timeout
	log.Println("Monitoring MongoDB connection for inactivity started.")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			service.mu.Lock()
			if time.Since(service.lastActivity) > service.timeout {
				log.Println("Closing MongoDB connection due to inactivity.")
				service.mu.Unlock() // Unlock before calling Close to avoid deadlock
				err := service.Close()
				if err != nil {
					log.Println("Error closing duo to inactivity MongoDB connection:", err)
				}
				return
			}
			service.mu.Unlock()
		case <-service.quit:
			return
		}
	}
}

func (service *MongoDBClient) getCollection(collection string) (*mongo.Collection, error) {
	// getCollection gets a MongoDB collection and updates the last activity time
	log.Println("Getting collection: ", collection)
	err := service.getClientConnection(URI)
	if err != nil {
		return nil, err
	}
	service.lastActivity = time.Now()
	return service.client.Database(DbName).Collection(collection), nil
}

func (service *MongoDBClient) Clear() error {
	// function that drops all records from the database
	log.Println("Clearing database")
	err := service.client.Database(DbName).Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Close manually closes the MongoDB connection
func (service *MongoDBClient) Close() error {
	if !service.isConnected() {
		return nil
	}
	log.Println("Closing MongoDB connection")
	if service.quit != nil {
		close(service.quit)
		service.quit = nil
	}
	return service.client.Disconnect(context.Background())
}
