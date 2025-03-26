package database

import (
	"app/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	// InsertAthlete inserts a new athlete into the database.
	InsertAthlete(athlete *model.StravaAthlete) error
	// GetAthleteIndex retrieves an athlete from the database by their name.
	GetAthleteIndex(athlete *model.StravaAthlete) error
	// GetAthleteByID retrieves an athlete from the database by their ID.
	GetAthleteByID(athleteId *primitive.ObjectID) (*model.StravaAthlete, error)
	// GetUniqueYears retrieves unique years from the activities collection.
	GetUniqueYears() ([]string, error)
	// GetLatestActivity retrieves the latest activity from the database.
	GetLatestActivity() (*model.StravaActivity, error)
	// InsertActivity inserts a new activity into the database.
	InsertActivity(activity *model.StravaActivity) error
	// UpdateAthleteDataSumUp updates the sum up of all activities for an athlete for a given year.
	UpdateAthleteDataSumUp(athleteId *primitive.ObjectID, year string) error
	// GetAthletesData retrieves the sum up of all activities for all athletes for a given year.
	GetAthletesData(year string) []model.AthleteData
	// RemoveActivities removes all activities from the database for a given date.
	RemoveActivities() error
}

func (s *service) InsertAthlete(athlete *model.StravaAthlete) error {
	log.Println("Inserting athlete")
	collection, err := s.getCollection(athletesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"firstname": athlete.Firstname, "lastname": athlete.Lastname}
	err = collection.FindOne(context.Background(), filter).Decode(athlete)
	if !errors.Is(err, mongo.ErrNoDocuments) {
		log.Fatalf("Athlete already exists: %v", err)
	}
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	if athlete.ID == nil {
		newID := primitive.NewObjectID()
		athlete.ID = &newID
	}

	_, err = collection.InsertOne(context.Background(), athlete)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAthleteIndex(athlete *model.StravaAthlete) error {
	log.Println(fmt.Sprintf("Getting athlete index: %s", athlete.ID))
	collection, err := s.getCollection(athletesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"firstname": athlete.Firstname, "lastname": athlete.Lastname}
	err = collection.FindOne(context.Background(), filter).Decode(&athlete)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = s.InsertAthlete(athlete)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAthleteByID(athleteId *primitive.ObjectID) (*model.StravaAthlete, error) {
	collection, err := s.getCollection(athletesCollection)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": athleteId}
	var athlete model.StravaAthlete
	err = collection.FindOne(context.Background(), filter).Decode(&athlete)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &athlete, nil
}

func (s *service) GetUniqueYears() ([]string, error) {
	log.Println("Getting unique years")
	collection, err := s.getCollection(activitiesCollection)
	if err != nil {
		return nil, err
	}

	// Use the distinct method to get unique date values
	dates, err := collection.Distinct(context.Background(), "date", bson.M{})
	if err != nil {
		return nil, err
	}

	uniqueYears := make(map[string]struct{})
	for _, date := range dates {
		// Extract the year part from the date string
		year := strings.Split(date.(string), "-")[0]
		uniqueYears[year] = struct{}{}
	}

	years := make([]string, 0, len(uniqueYears))
	for year := range uniqueYears {
		years = append(years, year)
	}
	return years, nil
}

func (s *service) GetLatestActivity() (*model.StravaActivity, error) {
	// official Strava API does not provide any ID for the activities,
	// so to avoid duplicates of the same activity in the database we need to get the latest activity
	log.Println("Getting latest activity")
	collection, err := s.getCollection(activitiesCollection)
	if err != nil {
		return nil, err
	}
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}}) // Sort by _id in descending order

	var activity model.StravaActivity
	err = collection.FindOne(context.Background(), bson.D{}, opts).Decode(&activity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &model.StravaActivity{}, nil
	} else if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *service) InsertActivity(activity *model.StravaActivity) error {
	log.Println(fmt.Sprintf("Inserting activity: %s, for athlete %s", activity.Name, activity.ID))
	collection, err := s.getCollection(activitiesCollection)
	if err != nil {
		return err
	}
	err = s.GetAthleteIndex(&activity.Athlete)
	if err != nil {
		return err
	}
	activity.UserID = activity.Athlete.ID

	if activity.ID == nil {
		newID := primitive.NewObjectID()
		activity.ID = &newID
	}

	if activity.Date == "" {
		activity.Date = time.Now().Format("2006-01-02")
	}

	_, err = collection.InsertOne(context.Background(), activity)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateAthleteDataSumUp(athleteId *primitive.ObjectID, year string) error {
	// TODO make it as transaction
	log.Println(fmt.Sprintf("Updateing sum up for %s, from: %s", athleteId, year))
	collection, err := s.getCollection(athleteDataSumCollection)
	if err != nil {
		log.Println("Error getting collection to  update athlete data sum up", err)
		return err
	}

	filter := bson.M{
		"userId": athleteId,
		"year":   year,
	}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Error deleting athlete data sum up", err)
		return err
	}
	athlete, err := s.GetAthleteByID(athleteId)
	if err != nil {
		log.Println("Error getting athlete by id to sum up", err)
		return err
	}
	athleteDataSum := s.calculateAthleteDataSumUp(athlete, year)
	err = s.insertAthleteDataSumUp(&athleteDataSum)
	if err != nil {
		log.Println("Error inserting athlete data sum up", err)
		return err
	}

	return nil
}

func (s *service) GetAthletesData(year string) []model.AthleteData {
	if year == "" {
		year = time.Now().Format("2006")
	}
	log.Println(fmt.Sprintf("Getting athletes data for year: %s", year))
	athleteData := make([]model.AthleteData, 0)
	athletes := s.getAthletes()

	for _, athlete := range athletes {
		athleteDataSumUp, _ := s.getAthleteDataSumUp(&athlete, year)
		if athleteDataSumUp.TotalActivities > 0 {
			athleteData = append(athleteData, athleteDataSumUp)
		}
	}
	return athleteData
}

func (s *service) RemoveActivities() error {
	date := "2025-02-27"
	log.Println(fmt.Sprintf("Removing all activities from date: %s", date))
	collection, err := s.getCollection(activitiesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"date": date}
	_, err = collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
