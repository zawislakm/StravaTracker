package database

import (
	"app/internal/model"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"math"
)

func (s *service) getAthletes() []model.StravaAthlete {
	log.Println("Getting all athletes")
	collection, err := s.getCollection(athletesCollection)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	athletes := make([]model.StravaAthlete, 0)
	if err := cursor.All(context.Background(), &athletes); err != nil {
		log.Fatal(err)
	}
	return athletes
}

func (s *service) getAthleteActivities(athlete *model.StravaAthlete, year string) []model.StravaActivity {
	log.Println(fmt.Sprintf("Getting all activities for athlete: %s", athlete.ID))
	collection, err := s.getCollection(activitiesCollection)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{
		"userId": athlete.ID,
		"date":   bson.M{"$regex": fmt.Sprintf("^%s", year)},
		//"type":   "Ride", // TODO investigate if there is a possibility to get other sport types from this club API
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	activities := make([]model.StravaActivity, 0)
	if err := cursor.All(context.Background(), &activities); err != nil {
		log.Fatal(err)
	}

	return activities
}

func (s *service) calculateAthleteDataSumUp(athlete *model.StravaAthlete, year string) model.AthleteData {
	log.Println(fmt.Sprintf("Getting sum up of all activities for athlete: %s", athlete.ID))

	activities := s.getAthleteActivities(athlete, year)

	athleteData := model.AthleteData{Name: athlete.Firstname + " " + athlete.Lastname, UserID: athlete.ID, Year: year}

	if len(activities) == 0 {
		return athleteData
	}

	athleteData.TotalActivities = len(activities)

	for _, activity := range activities {
		athleteData.Distance += activity.Distance
		athleteData.ElevationGain += activity.TotalElevationGain
		athleteData.LongestActivity = math.Max(athleteData.LongestActivity, activity.Distance)
		athleteData.TotalTime += float64(activity.MovingTime)
	}

	// convert distance from meters to kilometers
	athleteData.Distance /= 1000
	athleteData.LongestActivity /= 1000

	athleteData.AverageTime = athleteData.TotalTime / float64(athleteData.TotalActivities)
	athleteData.AverageLength = athleteData.Distance / float64(athleteData.TotalActivities)
	athleteData.AverageSpeed = athleteData.Distance / (athleteData.TotalTime / 3600)
	return athleteData
}

func (s *service) insertAthleteDataSumUp(athleteData *model.AthleteData) error {
	log.Println(fmt.Sprintf("Inserting sum up of athlete: %s, for athlete: %s", athleteData, athleteData.UserID))
	collection, err := s.getCollection(athleteDataSumCollection)
	if err != nil {
		log.Println("Error getting collection to  insert athlete data sum up", err)
		return err
	}

	_, err = collection.InsertOne(context.Background(), athleteData)
	if err != nil {
		log.Println("Error inserting athlete data sum up", err)
		return err
	}
	return nil
}

func (s *service) getAthleteDataSumUp(athlete *model.StravaAthlete, year string) (model.AthleteData, error) {
	log.Println(fmt.Sprintf("Getting sum up of all activities for athlete: %s from year: %s", athlete.ID, year))
	collection, err := s.getCollection(athleteDataSumCollection)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{
		"userId": athlete.ID,
		"year":   year,
	}

	var athleteDataSum model.AthleteData
	err = collection.FindOne(context.Background(), filter).Decode(&athleteDataSum)

	if errors.Is(err, mongo.ErrNoDocuments) {
		athleteDataSum = s.calculateAthleteDataSumUp(athlete, year)
		err = s.insertAthleteDataSumUp(&athleteDataSum)
		if err != nil {
			log.Println(fmt.Sprintf("Error inserting athlete data: %v", err))
			return athleteDataSum, err
		}
	} else if err != nil {
		return athleteDataSum, err
	}

	return athleteDataSum, nil
}
