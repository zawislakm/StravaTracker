package server

import (
	"app/internal/database"
	"log"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"app/internal/strava"
	_ "github.com/joho/godotenv/autoload"
)

var apiCallTimeout = 1 * time.Minute

func GetActivities(apiService *strava.ServiceStravaAPI, dbService database.Service, newDataChan chan bool) {
	// fill this function to run each 5 minutes
	// get the latest activities from the Strava API
	log.Println("Goroutine to get activities started")
	ticker := time.NewTicker(apiCallTimeout)

	for {
		select {
		case <-ticker.C:
			// get the latest activities from the Strava API
			log.Println("Calling for activities")
			activities, err := apiService.StravaGetClubActivities()
			if err != nil {
				log.Printf("Error getting activities from Strava API: %v", err)
			} else {
				newActivities := filterNewActivities(activities, dbService)

				if len(newActivities) > 0 {
					if len(newActivities) > 1000 { // TODO, find out why it happens, I guess its error from strava
						log.Printf("New acctiviteid found, more than 10: %v, something wrong happend, skipping this request data \n", len(newActivities))
					} else {
						processNewActivities(newActivities, dbService)
						log.Println("New activities found")
						// TODO send a notification to frontend about new activities
						newDataChan <- true // pass year of activities to reload, maybe reload is not needed
					}
				} else {
					log.Println("No new activities found")
				}
			}
		}

	}
}

type Server struct {
	port   int
	strava *strava.ServiceStravaAPI
	db     database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		strava: strava.GetStravaClient(),
		db:     database.GetDbClient(),
	}
	newDataChan := make(chan bool)
	go waitForNewData(newDataChan)
	go GetActivities(NewServer.strava, NewServer.db, newDataChan)

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
