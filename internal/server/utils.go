package server

import (
	"app/internal/database"
	"app/internal/model"
	"fmt"
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

type athleteUpdateKey struct {
	id   *primitive.ObjectID
	year string
}

func (auk athleteUpdateKey) keyValue() string {
	return auk.id.Hex() + "_" + auk.year
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

	athletesRequiringUpdate := map[string]athleteUpdateKey{}
	for i := len(activities) - 1; i >= 0; i-- {
		err := dbService.InsertActivity(&activities[i]) //pass pointer so it would be updated
		if err != nil {
			log.Printf("Error inserting activity: %v\n", err)
			continue
		}
		key := athleteUpdateKey{
			id:   activities[i].UserID,
			year: getYear(activities[i].Date),
		}
		fmt.Println(key, key.id, key.year, activities[i].Date, getYear(activities[i].Date), activities[i].Date)
		athletesRequiringUpdate[key.keyValue()] = key
	}
	fmt.Println(len(activities), len(athletesRequiringUpdate))
	for _, athlete := range athletesRequiringUpdate {
		fmt.Println(athlete, athlete.id, athlete.year)
		err := dbService.UpdateAthleteDataSumUp(athlete.id, athlete.year)
		if err != nil {
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
