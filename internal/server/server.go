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

var apiCallTimeout, _ = strconv.Atoi(os.Getenv("API_CALL_TIMEOUT"))

func (s *Server) getActivities() {
	// get the latest activities from the Strava API
	log.Println("Goroutine to get activities started")

	for {
		// get the latest activities from the Strava API
		log.Println("Calling for activities")
		activities, err := s.strava.StravaGetClubActivities()

		if err != nil {
			log.Printf("Error getting activities from Strava API: %v", err)
		} else {
			s.lastUpdate = time.Now()
			log.Println("Called for activities, new lastUpdate:", s.lastUpdate)
			newActivities := s.filterNewActivities(activities)
			newActivitiesLength := len(newActivities)
			isNewActivities := false

			if newActivitiesLength > 100 {
				// TODO, find out why it happens, I guess its error from strava
				log.Printf("New acctiviteid found, more than 10: %v, something wrong happend, skipping this request data \n", len(newActivities))
			} else if newActivitiesLength > 0 {
				s.processNewActivities(newActivities)
				isNewActivities = true
				log.Println("New activities found")
			} else {
				log.Println("No activities found")
			}
			select {
			case s.newActivitiesChan <- isNewActivities:
			default:
			}
		}
		time.Sleep(time.Duration(apiCallTimeout) * time.Minute)
	}
}

type Server struct {
	port              int
	strava            strava.ServiceStravaAPI
	db                database.Service
	newActivitiesChan chan bool
	lastUpdate        time.Time
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:              port,
		strava:            strava.GetStravaClient(),
		db:                database.GetDbClient(),
		newActivitiesChan: make(chan bool),
	}
	go NewServer.getActivities()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 0, // Zero disables the write timeout, used in SSE
	} // possible to create a separate Server on other port only for SSE connection

	return server
}
