package Database

import (
	"app/src/Models"
	"context"
	"errors"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (service *MongoDBClient) insertAthlete(athlete *Models.StravaAthlete) error {
	collection, err := service.getCollection(dbVariables.AthletesCollection)
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

func (service *MongoDBClient) GetAthleteIndex(athlete *Models.StravaAthlete) error {
	collection, err := service.getCollection(dbVariables.AthletesCollection)
	if err != nil {
		return err
	}

	filter := bson.M{"firstname": athlete.Firstname, "lastname": athlete.Lastname}
	err = collection.FindOne(context.Background(), filter).Decode(&athlete)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = service.insertAthlete(athlete)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (service *MongoDBClient) GetLatestActivity() (*Models.StravaActivity, error) {
	// official Strava API does not provide any ID for the activities,
	// so to avoid duplicates of the same activity in the database we need to get the latest activity
	collection, err := service.getCollection(dbVariables.ActivitiesCollection)
	if err != nil {
		return nil, err
	}
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}}) // Sort by _id in descending order

	var activity Models.StravaActivity
	err = collection.FindOne(context.Background(), bson.D{}, opts).Decode(&activity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &Models.StravaActivity{}, nil
	} else if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (service *MongoDBClient) InsertActivity(activity Models.StravaActivity) error {
	collection, err := service.getCollection(dbVariables.ActivitiesCollection)
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
	return nil
}

func (service *MongoDBClient) getAthletes() []Models.StravaAthlete {
	collection, err := service.getCollection(dbVariables.AthletesCollection)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	athletes := make([]Models.StravaAthlete, 0)
	if err := cursor.All(context.Background(), &athletes); err != nil {
		log.Fatal(err)
	}
	return athletes
}

func (service *MongoDBClient) getAthleteActivities(athlete *Models.StravaAthlete) []Models.StravaActivity {
	collection, err := service.getCollection(dbVariables.ActivitiesCollection)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"userId": athlete.ID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	activities := make([]Models.StravaActivity, 0)
	if err := cursor.All(context.Background(), &activities); err != nil {
		log.Fatal(err)
	}

	return activities
}

func (service *MongoDBClient) getAthleteDataSumUp(athlete *Models.StravaAthlete) Models.AthleteData {
	activities := service.getAthleteActivities(athlete)

	athleteData := Models.AthleteData{Name: athlete.Firstname + " " + athlete.Lastname}

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

func (service *MongoDBClient) GetAthletesData() []Models.AthleteData {
	athleteData := make([]Models.AthleteData, 0)
	athletes := service.getAthletes()

	for _, athlete := range athletes {
		athleteData = append(athleteData, service.getAthleteDataSumUp(&athlete))
	}

	return athleteData
}
