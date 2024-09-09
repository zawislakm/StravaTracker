package Server

import (
	"app/src/Database"
	"app/src/Models"
	"app/src/StravaAPI"
	"log"
	"time"
)

var apiCallTimeout = 10 * time.Minute

func GetActivities(apiService *StravaAPI.ServiceStravaAPI, dbService *Database.MongoDBClient) {
	// fill this function to run each 5 minutes
	// get the latest activities from the Strava API
	log.Println("Go routine to get activities started")
	ticker := time.NewTicker(apiCallTimeout)

	for {
		select {
		case <-ticker.C:
			// get the latest activities from the Strava API
			log.Println("Calling for activities")
			activities, err := apiService.StravaGetClubActivities()
			if err != nil {
				// log the error
			}
			newActivities := filterNewActivities(activities, dbService)
			processNewActivities(newActivities, dbService)
			if len(newActivities) > 0 {
				log.Println("New activities found")
				// TODO send a notification to frontend about new activities
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
