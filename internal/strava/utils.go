package strava

import (
	"app/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type serviceStravaAPI struct {
	Client      *http.Client
	CachedToken *model.StravaOauthResponse
}

type apiStravaVariables struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	ClubID       string
}

var (
	stravaURL      = "https://www.strava.com/api/v3"
	defaultPerPage = 30

	apiClient    *serviceStravaAPI
	apiVariables *apiStravaVariables
)

func init() {
	apiClient = &serviceStravaAPI{
		Client:      &http.Client{},
		CachedToken: nil,
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	apiVariables = &apiStravaVariables{
		ClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		RefreshToken: os.Getenv("STRAVA_REFRESH_TOKEN"),
		ClubID:       os.Getenv("STRAVA_CLUB_ID"),
	}
	log.Println("StravaAPI variables initialized")
}

func GetStravaClient() ServiceStravaAPI {
	log.Println("Getting Strava client service")
	if apiClient == nil {
		apiClient = &serviceStravaAPI{
			Client:      &http.Client{},
			CachedToken: nil,
		}
	}
	return apiClient
}

func stravaGetToken() (model.StravaOauthResponse, error) {

	baseURL := fmt.Sprintf("%s/oauth/token", stravaURL)

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
		if err := Body.Close(); err != nil {

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

func getStravaHeader(service *serviceStravaAPI) (model.StravaHeader, error) {
	if service.CachedToken == nil || service.CachedToken.IsExpired() {
		oauthResponse, err := stravaGetToken()
		if err != nil {
			return model.StravaHeader{}, err
		}
		service.CachedToken = &oauthResponse
	}

	return model.StravaHeader{
		Authorization: fmt.Sprintf("Bearer %s", service.CachedToken.AccessToken),
	}, nil
}

func (s *serviceStravaAPI) setRequestHeaders(req *http.Request) error {
	headers, err := getStravaHeader(s)
	if err != nil {
		return fmt.Errorf("error getting headers: %v", err)
	}

	req.Header.Add("Authorization", headers.Authorization)
	req.Header.Add("Content-Type", "application/json")
	return nil
}

func (s *serviceStravaAPI) setUpRequest(method string, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %v", path, err)
	}

	if err = s.setRequestHeaders(req); err != nil {
		return nil, fmt.Errorf("error setting headers for %s: %v", path, err)
	}

	return req, nil
}

func performRequest[T model.StravaAthlete | model.StravaActivity](s *serviceStravaAPI, method string, path string, responseSave []T) ([]T, error) {
	req, err := s.setUpRequest(method, path)
	if err != nil {
		return nil, fmt.Errorf("error setting up request for %s: %v", path, err)
	}

	response, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request for %s: %v", path, err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status for %s: %s", path, response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseSave); err != nil {
		return nil, fmt.Errorf("error decoding response for %s: %v", path, err)
	}
	return responseSave, nil
}

func iteratePages[T model.StravaAthlete | model.StravaActivity](s *serviceStravaAPI, method string, path string, responseSave []T) ([]T, error) {
	var responseData []T
	var page = 1
	for {
		params := fmt.Sprintf("?page=%d&per_page=%d", page, defaultPerPage)
		fullURL := fmt.Sprintf("%s%s", path, params)

		elem, err := performRequest(s, method, fullURL, responseSave)
		if err != nil {
			return nil, err
		}
		if len(elem) == 0 {
			break
		}
		responseData = append(responseData, elem...)
		page++
	}
	return responseData, nil
}
