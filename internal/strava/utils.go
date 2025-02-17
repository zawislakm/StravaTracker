package strava

import (
	"app/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var StravaURL = "https://www.strava.com/api/v3"
var DefaultPerPage = 30

var apiClient *ServiceStravaAPI
var apiVariables *ApiStravaVariables

type ServiceStravaAPI struct {
	Client      *http.Client
	CachedToken *model.StravaOauthResponse
}

type ApiStravaVariables struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	ClubID       string
}

func init() {
	apiClient = &ServiceStravaAPI{
		Client:      &http.Client{},
		CachedToken: nil,
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	apiVariables = &ApiStravaVariables{
		ClientID:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
		RefreshToken: os.Getenv("STRAVA_REFRESH_TOKEN"),
		ClubID:       os.Getenv("STRAVA_CLUB_ID"),
	}
	log.Println("Strava API initialized")
}

func GetStravaClient() *ServiceStravaAPI {
	if apiClient == nil {
		apiClient = &ServiceStravaAPI{
			Client:      &http.Client{},
			CachedToken: nil,
		}
	}
	return apiClient
}

func getStravaHeader(service *ServiceStravaAPI) (model.StravaHeader, error) {
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

func (service *ServiceStravaAPI) setRequestHeaders(req *http.Request) error {
	headers, err := getStravaHeader(service)
	if err != nil {
		return fmt.Errorf("Error getting headers: %v", err)
	}

	req.Header.Add("Authorization", headers.Authorization)
	req.Header.Add("Content-Type", "application/json")
	return nil
}

func (service *ServiceStravaAPI) setUpRequest(method string, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %v", path, err)
	}

	err = service.setRequestHeaders(req)
	if err != nil {
		return nil, fmt.Errorf("error setting headers for %s: %v", path, err)
	}

	return req, nil
}

func performRequest[T model.StravaAthlete | model.StravaActivity](service *ServiceStravaAPI, method string, path string, responseSave []T) ([]T, error) {
	req, err := service.setUpRequest(method, path)
	if err != nil {
		return nil, fmt.Errorf("error setting up request for %s: %v", path, err)
	}

	response, err := service.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request for %s: %v", path, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
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

func iteratePages[T model.StravaAthlete | model.StravaActivity](service *ServiceStravaAPI, method string, path string, responseSave []T) ([]T, error) {
	var responseData []T
	var page = 1
	for {
		params := fmt.Sprintf("?page=%d&per_page=%d", page, DefaultPerPage)
		fullURL := fmt.Sprintf("%s%s", path, params)

		elem, err := performRequest(service, method, fullURL, responseSave)
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
