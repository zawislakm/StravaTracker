package Models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func (s *StravaOauthResponse) IsExpired() bool {
	return time.Now().Unix() > int64(s.ExpiresAt)
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

func (apiActivity *StravaActivity) CompareStravaData(dbActivity *StravaActivity) bool {
	return apiActivity.Name == dbActivity.Name && apiActivity.Distance == dbActivity.Distance && apiActivity.MovingTime == dbActivity.MovingTime && apiActivity.TotalElevationGain == dbActivity.TotalElevationGain && apiActivity.Type == dbActivity.Type && apiActivity.SportType == dbActivity.SportType && apiActivity.WorkoutType == dbActivity.WorkoutType
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

func (athleteData *AthleteData) FormatDuration() string {
	hours := int(athleteData.AverageTime) / 3600
	minutes := (int(athleteData.AverageTime) % 3600) / 60
	seconds := int(athleteData.AverageTime) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh:%02dm:%02ds", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02dm:%02ds", minutes, seconds)
}
