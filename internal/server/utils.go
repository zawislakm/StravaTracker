package server

import (
	"app/cmd/web/templates"
	"app/internal/model"
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func (s *Server) filterNewActivities(activities []model.StravaActivity) []model.StravaActivity {
	// get the latest activity from the database
	// filter the activities that are not in the database

	if len(activities) == 0 {
		return []model.StravaActivity{}
	}

	mostRecentDBActivity, err := s.db.GetLatestActivity()
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

func (s *Server) processNewActivities(activities []model.StravaActivity) {
	// insert the new activities into the database in reverse order to the newest activity is inserted last

	athletesRequiringUpdate := map[string]athleteUpdate{}
	for i := len(activities) - 1; i >= 0; i-- { // adding in reverse orders so the newest activity was added last to db
		activity := activities[i]
		if err := s.db.InsertActivity(&activity); err != nil {
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
		if err := s.db.UpdateAthleteDataSumUp(athlete.id, athlete.year); err != nil {
			log.Printf("Error updating athlete data sum up: %v\n", err)
		}
	}
}

func sendDateUpdate(c echo.Context, lastUpdate time.Time) {
	var buf bytes.Buffer
	err := templates.Update(lastUpdate.Format("2006-01-02 15:04")).Render(c.Request().Context(), &buf)
	if err != nil {
		fmt.Printf("Error rendering update template: %v\n", err)
		return
	}
	html := buf.String()
	fmt.Fprintf(c.Response(), "event: Update\ndata: %s\n\n", html)
	c.Response().Flush()
}
