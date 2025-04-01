package strava

import (
	"app/internal/model"
	"fmt"
	"net/http"
)

type ServiceStravaAPI interface {
	// StravaGetClubAthletes retrieves the athletes in the specified Strava club.
	StravaGetClubAthletes() ([]model.StravaAthlete, error)
	// StravaGetClubActivities retrieves the activities of the specified Strava club.
	StravaGetClubActivities() ([]model.StravaActivity, error)
}

func (s *serviceStravaAPI) StravaGetClubAthletes() ([]model.StravaAthlete, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/members", stravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(s, http.MethodGet, baseURL, []model.StravaAthlete{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func (s *serviceStravaAPI) StravaGetClubActivities() ([]model.StravaActivity, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/activities", stravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(s, http.MethodGet, baseURL, []model.StravaActivity{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
