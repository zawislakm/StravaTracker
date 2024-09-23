package StravaAPI

import (
	"app/src/Models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func stravaGetToken() (Models.StravaOauthResponse, error) {

	baseURL := fmt.Sprintf("%s/oauth/token", StravaURL)

	payload := Models.StravaOauthRequest{
		ClientID:     apiVariables.ClientID,
		ClientSecret: apiVariables.ClientSecret,
		RefreshToken: apiVariables.RefreshToken,
		GrantType:    "refresh_token",
		F:            "json",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Models.StravaOauthResponse{}, err
	}

	response, err := http.Post(baseURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil {
		return Models.StravaOauthResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return Models.StravaOauthResponse{}, fmt.Errorf("unexpected status: %s", response.Status)
	}

	var responseData Models.StravaOauthResponse
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return Models.StravaOauthResponse{}, err
	}

	return responseData, nil
}

func (service *ServiceStravaAPI) StravaGetClubAthletes() ([]Models.StravaAthlete, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/members", StravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(service, http.MethodGet, baseURL, []Models.StravaAthlete{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}

func (service *ServiceStravaAPI) StravaGetClubActivities() ([]Models.StravaActivity, error) {
	baseURL := fmt.Sprintf("%s/clubs/%s/activities", StravaURL, apiVariables.ClubID)
	responseData, err := iteratePages(service, http.MethodGet, baseURL, []Models.StravaActivity{})
	if err != nil {
		return nil, err
	}
	return responseData, nil
}
