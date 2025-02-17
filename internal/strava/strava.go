package strava

import (
	"app/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func stravaGetToken() (model.StravaOauthResponse, error) {

	baseURL := fmt.Sprintf("%s/oauth/token", StravaURL)

	payload := model.StravaOauthRequest{
		ClientID:     apiVariables.ClientID,
		ClientSecret: apiVariables.ClientSecret,
		RefreshToken: apiVariables.RefreshToken,
		GrantType:    "refresh_token",
		F:            "json",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return model.StravaOauthResponse{}, err
	}

	response, err := http.Post(baseURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil {
		return model.StravaOauthResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return model.StravaOauthResponse{}, fmt.Errorf("unexpected status: %s", response.Status)
	}

	var responseData model.StravaOauthResponse
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return model.StravaOauthResponse{}, err
	}

	return responseData, nil
}

func (service *ServiceStravaAPI) StravaGetClubAthletes() ([]model.StravaAthlete, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/members", StravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(service, http.MethodGet, baseURL, []model.StravaAthlete{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func (service *ServiceStravaAPI) StravaGetClubActivities() ([]model.StravaActivity, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/activities", StravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(service, http.MethodGet, baseURL, []model.StravaActivity{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
