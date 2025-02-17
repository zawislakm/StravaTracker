package server

import (
	"app/internal/database"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"app/internal/strava"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port   int
	strava *strava.ServiceStravaAPI
	db     *database.MongoDBClient
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		strava: strava.GetStravaClient(),
		db:     database.GetDbClient(),
	}
	//cache := strava.NewDataCache(NewServer.db)
	//go strava.GetActivities(NewServer.strava, NewServer.db, cache)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
