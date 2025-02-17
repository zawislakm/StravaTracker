package database

import (
	"app/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO rethink service locking logic, maybe it is not needed

func (service *MongoDBClient) InsertAthlete(athlete *model.StravaAthlete) error {
	//service.mu.Lock()
	log.Println("Inserting athlete")
	collection, err := service.getCollection(athletesCollection)
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
	////defer service.mu.Unlock()
	return nil
}

func (service *MongoDBClient) GetAthleteIndex(athlete *model.StravaAthlete) error {
	//service.mu.Lock()
	////defer service.mu.Unlock()
	log.Println(fmt.Sprintf("Getting athlete index: %s", athlete.ID))
	collection, err := service.getCollection(athletesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"firstname": athlete.Firstname, "lastname": athlete.Lastname}
	err = collection.FindOne(context.Background(), filter).Decode(&athlete)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = service.InsertAthlete(athlete)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	//service.mu.Unlock()
	return nil
}

func (service *MongoDBClient) GetUniqueYears() ([]string, error) {
	//service.mu.Lock()
	log.Println("Getting unique years")
	collection, err := service.getCollection(activitiesCollection)
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
	//service.mu.Unlock()
	return years, nil
}

func (service *MongoDBClient) GetLatestActivity() (*model.StravaActivity, error) {
	// official Strava API does not provide any ID for the activities,
	// so to avoid duplicates of the same activity in the database we need to get the latest activity
	//service.mu.Lock()
	log.Println("Getting latest activity")
	collection, err := service.getCollection(activitiesCollection)
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
	//service.mu.Unlock()
	return &activity, nil
}

func (service *MongoDBClient) InsertActivity(activity model.StravaActivity) error {
	//service.mu.Lock()
	log.Println(fmt.Sprintf("Inserting activity: %s, for athlete %s", activity.Name, activity.ID))
	collection, err := service.getCollection(activitiesCollection)
	if err != nil {
		return err
	}
	err = service.GetAthleteIndex(&activity.Athlete)
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
	//service.mu.Unlock()
	return nil
}

func (service *MongoDBClient) getAthletes() []model.StravaAthlete {
	log.Println("Getting all athletes")
	collection, err := service.getCollection(athletesCollection)
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

func (service *MongoDBClient) getAthleteActivities(athlete *model.StravaAthlete, year string) []model.StravaActivity {
	log.Println(fmt.Sprintf("Getting all activities for athlete: %s", athlete.ID))
	collection, err := service.getCollection(activitiesCollection)
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

func (service *MongoDBClient) getAthleteDataSumUp(athlete *model.StravaAthlete, year string) model.AthleteData {
	log.Println(fmt.Sprintf("Getting sum up of all activities for athlete: %s", athlete.ID))
	activities := service.getAthleteActivities(athlete, year)

	athleteData := model.AthleteData{Name: athlete.Firstname + " " + athlete.Lastname}

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

func (service *MongoDBClient) GetAthletesData(year string) []model.AthleteData {
	//service.mu.Lock()
	log.Println("Getting sum up of all activities for all athletes")
	if year == "" {
		year = time.Now().Format("2006")
	}
	log.Println(fmt.Sprintf("Year: %s", year))
	athleteData := make([]model.AthleteData, 0)
	athletes := service.getAthletes()

	for _, athlete := range athletes {
		athleteDataSumUp := service.getAthleteDataSumUp(&athlete, year)
		if athleteDataSumUp.TotalActivities > 0 {
			athleteData = append(athleteData, athleteDataSumUp)
		}
	}
	//service.mu.Unlock()
	return athleteData
}

// remove all activities from given date
func (service *MongoDBClient) RemoveActivities() error {
	date := "2025-02-27"
	//service.mu.Lock()
	log.Println(fmt.Sprintf("Removing all activities from date: %s", date))
	collection, err := service.getCollection(activitiesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"date": date}
	_, err = collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	//service.mu.Unlock()
	return nil
}
