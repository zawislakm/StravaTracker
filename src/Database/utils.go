package Database

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

var once sync.Once
var dbClient *MongoDBClient
var dbVariables *MongoDBVariables

//TODO add logging
//TODO better error handling
//TODO add tests
//TODO add docstrings

type MongoDBClient struct {
	client       *mongo.Client
	lastActivity time.Time
	timeout      time.Duration
	mu           sync.Mutex
	quit         chan struct{}
}

type MongoDBVariables struct {
	URI                  string
	Database             string
	AthletesCollection   string
	ActivitiesCollection string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbVariables = &MongoDBVariables{
		URI:                  os.Getenv("DB_URI"),
		Database:             os.Getenv("DB_NAME"),
		AthletesCollection:   os.Getenv("DB_ATHLETES_COLLECTION"),
		ActivitiesCollection: os.Getenv("DB_ACTIVITIES_COLLECTION"),
	}

	dbClient = &MongoDBClient{
		client:       nil,
		lastActivity: time.Now(),
		timeout:      defaultTimeout,
		mu:           sync.Mutex{},
		quit:         nil,
	}

	err := dbClient.getClientConnection(dbVariables.URI)
	if err != nil {
		log.Fatal(err)
	}

	once.Do(func() {
		// TODO check if setting indexes is necessary
		err := dbClient.setAthleteIndex()
		if err != nil {
			log.Fatal(err)
		}
		//err = dbClient.setActivityIndex()
		//if err != nil {
		//	log.Fatal(err)
		//}
	})

	fmt.Println("Database package initialized")
}

func GetDbClient() *MongoDBClient {
	if dbClient == nil {
		dbClient = &MongoDBClient{
			client:       nil,
			lastActivity: time.Now(),
			timeout:      defaultTimeout,
			mu:           sync.Mutex{},
			quit:         nil,
		}
	}
	return dbClient
}

func (service *MongoDBClient) setAthleteIndex() error {
	collection, err := service.getCollection(dbVariables.AthletesCollection)
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
	collection, err := service.getCollection(dbVariables.ActivitiesCollection)
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

// getClientConnection creates a new MongoDB client with an inactivity timeout
func (service *MongoDBClient) getClientConnection(uri string) error {

	if service.isConnected() {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(dbVariables.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	service.client = client
	service.quit = make(chan struct{})
	go service.monitorInactivity()

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

// monitorInactivity monitors the client's activity and closes the connection after timeout
func (service *MongoDBClient) monitorInactivity() {
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

// getCollection gets a MongoDB collection and updates the last activity time
func (service *MongoDBClient) getCollection(collection string) (*mongo.Collection, error) {
	err := service.getClientConnection(dbVariables.URI)
	if err != nil {
		return nil, err
	}
	service.mu.Lock()
	defer service.mu.Unlock()

	service.lastActivity = time.Now()
	return service.client.Database(dbVariables.Database).Collection(collection), nil
}

// Close manually closes the MongoDB connection
func (service *MongoDBClient) Close() error {
	if !service.isConnected() {
		return nil
	}
	if service.quit != nil {
		close(service.quit)
		service.quit = nil
	}
	return service.client.Disconnect(context.Background())
}
