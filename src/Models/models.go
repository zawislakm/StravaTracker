package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StravaOauthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
	F            string `json:"f"`
}

type StravaOauthResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type StravaHeader struct {
	Authorization string `json:"Authorization"`
}

type StravaAthlete struct {
	ID            *primitive.ObjectID `bson:"_id"`
	ResourceState int                 `json:"resource_state" bson:"-"`
	Firstname     string              `json:"firstname" bson:"firstname"`
	Lastname      string              `json:"lastname" bson:"lastname"`
	Membership    string              `json:"membership" bson:"-"`
	Admin         bool                `json:"admin" bson:"-"`
	Owner         bool                `json:"owner" bson:"-"`
}

type StravaActivity struct {
	ID                 *primitive.ObjectID `bson:"_id"`
	UserID             *primitive.ObjectID `bson:"userId"`
	ResourceState      int                 `json:"resource_state" bson:"-"`
	Athlete            StravaAthlete       `json:"athlete" bson:"-"`
	Name               string              `json:"name" bson:"name"`
	Distance           float64             `json:"distance" bson:"distance"`
	MovingTime         int                 `json:"moving_time" bson:"movingTime"`
	ElapsedTime        int                 `json:"elapsed_time" bson:"-"`
	TotalElevationGain float64             `json:"total_elevation_gain" bson:"totalElevationGain"`
	Type               string              `json:"type"  bson:"type"`
	SportType          string              `json:"sport_type" bson:"sportType"`
	WorkoutType        int                 `json:"workout_type" bson:"workoutType"`
	Date               string              `bson:"date"`
}

type AthleteData struct {
	Name            string
	Distance        float64
	TotalTime       float64
	LongestActivity float64
	ElevationGain   float64
	TotalActivities int
	AverageTime     float64
	AverageSpeed    float64
	AverageLength   float64
}
