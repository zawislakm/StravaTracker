package Server

import (
	"app/src/Database"
	"app/src/Models"
	"app/src/StravaAPI"
	"log"
	"sync"
	"time"
)

type DataCache struct {
	// idk if cache is needed now when year is stored here also, may lead to often reloading, rethink this
	Activities []Models.AthleteData
	Year       string
	ReloadChan chan bool
	mu         sync.Mutex
}

var apiCallTimeout = 5 * time.Minute

func GetActivities(apiService *StravaAPI.ServiceStravaAPI, dbService *Database.MongoDBClient, cache *DataCache) {
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
				log.Fatalf("Error getting activities from Strava API: %v", err)
			}
			newActivities := filterNewActivities(activities, dbService)
			processNewActivities(newActivities, dbService)
			if len(newActivities) > 0 {
				log.Println("New activities found")
				// TODO send a notification to frontend about new activities
				cache.ReloadChan <- true // pass year of activities to reload, maybe reload is not needed
			} else {
				log.Println("No new activities found")
			}
		}

	}
}

func filterNewActivities(activities []Models.StravaActivity, dbService *Database.MongoDBClient) []Models.StravaActivity {
	// get the latest activity from the database
	// filter the activities that are not in the database

	mostRecentDBActivity, err := dbService.GetLatestActivity()
	if err != nil {
		log.Fatalf("Error getting latest activity from database: %v", err)
	}

	newActivities := make([]Models.StravaActivity, 0)

	for _, activity := range activities {
		if activity.CompareStravaData(mostRecentDBActivity) {
			break
		}
		newActivities = append(newActivities, activity)
	}

	return newActivities
}

func processNewActivities(activities []Models.StravaActivity, dbService *Database.MongoDBClient) {
	// insert the new activities into the database in reverse order to the newest activity is inserted last
	for i := len(activities) - 1; i >= 0; i-- {
		err := dbService.InsertActivity(activities[i])
		if err != nil {
			log.Fatalf("Error inserting activity: %v", err)
		}
	}

}

func newDataCache(serviceDb *Database.MongoDBClient) *DataCache {
	log.Println("Creating data cache")
	year := time.Now().Format("2006")
	cache := &DataCache{
		Activities: serviceDb.GetAthletesData(year),
		Year:       year,
		ReloadChan: make(chan bool),
	}
	go cache.reloadData(serviceDb)
	return cache
}

func (cache *DataCache) GetActivities(serviceDb *Database.MongoDBClient, year string) []Models.AthleteData {
	log.Println("Getting activities from cache")
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.Year != year {
		cache.Year = year
		cache.Activities = serviceDb.GetAthletesData(year)
	}

	return cache.Activities
}

func (cache *DataCache) reloadData(serviceDb *Database.MongoDBClient) {
	for {
		select {
		case <-cache.ReloadChan:
			log.Println("Reloading cached athlete data")
			cache.mu.Lock()
			cache.Activities = serviceDb.GetAthletesData(cache.Year)
			cache.mu.Unlock()
		}
	}
}
