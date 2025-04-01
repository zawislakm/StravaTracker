package server

import (
	"app/internal/database"
	"app/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func filterNewActivities(activities []model.StravaActivity, dbService database.Service) []model.StravaActivity {
	// get the latest activity from the database
	// filter the activities that are not in the database

	mostRecentDBActivity, err := dbService.GetLatestActivity()
	if err != nil {
		// could not get the latest activity from the database, impossible to compare the activities to find new ones
		log.Printf("Error getting latest activity from database: %v \n", err)
		return []model.StravaActivity{}
	}

	newActivities := make([]model.StravaActivity, 0)

	for _, activity := range activities {
		if activity.CompareStravaData(mostRecentDBActivity) {
			break
		}
		newActivities = append(newActivities, activity)
	}

	return newActivities
}

type athleteUpdate struct {
	id   *primitive.ObjectID
	year string
}

func (au athleteUpdate) getKey() string {
	return au.id.Hex() + "_" + au.year
}

func getYear(dateStr string) string {
	data, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return ""
	}
	return data.Format("2006")
}

func processNewActivities(activities []model.StravaActivity, dbService database.Service) {
	// insert the new activities into the database in reverse order to the newest activity is inserted last

	athletesRequiringUpdate := map[string]athleteUpdate{}
	for _, activity := range activities {
		if err := dbService.InsertActivity(&activity); err != nil {
			log.Printf("Error inserting activity: %v\n", err)
			continue
		}
		au := athleteUpdate{
			id:   activity.UserID,
			year: getYear(activity.Date),
		}
		athletesRequiringUpdate[au.getKey()] = au
	}
	for _, athlete := range athletesRequiringUpdate {
		if err := dbService.UpdateAthleteDataSumUp(athlete.id, athlete.year); err != nil {
			log.Printf("Error updating athlete data sum up: %v\n", err)
		}
	}
}

func waitForNewData(newDataChan chan bool) {
	log.Println("Waiting for new data started")
	for {
		select {
		case <-newDataChan:
			log.Println("New data found chan function")
		}
	}
}
