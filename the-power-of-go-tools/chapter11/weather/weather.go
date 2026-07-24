package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type OWMDataAPIResponse struct {
	Name  string `json:"name"`
	Coord struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"coord"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
}

type OWMGeoAPIResponse struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Conditions struct {
	Summary string
}

const BaseURL = "https://api.openweathermap.org"
const DataAPI = "/data/2.5/weather?"
const GeoAPI = "/geo/1.0/direct?"
const BrisbaneLocation = "Brisbane,AU"
const BrisbaneLatitude = "-27.4651"
const BrisbaneLongitude = "153.0231"

// Returns the API key associated with environment variable OWM_API_KEY.
func APIKey() (string, error) {
	key := os.Getenv("OWM_API_KEY")
	if key == "" {
		msg := "error: environment variable OWM_API_KEY is not set"
		return "", errors.New(msg)
	}
	return key, nil
}

// Reads a []byte as JSON in the OWNDataResponse, makes sure there is at least 1
// Weather element, and converts OWNDataResponse into Conditions.
func ParseDataAPIResponse(data []byte) (Conditions, error) {
	var response OWMDataAPIResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return Conditions{}, fmt.Errorf("invalid API response %s: %w", data, err)
	}
	if len(response.Weather) > 1 {
		return Conditions{}, fmt.Errorf("invalid API response %s needs one Weather element", data)
	}
	conditions := Conditions{
		Summary: response.Weather[0].Main,
	}
	return conditions, nil
}

// Returns a URL using url.Values.Encode for encoding query parameters.
// Requires API endpoint to call, and the args map containing the key.
func FormatURL(api string, args map[string]string) string {
	// https://api.openweathermap.org/data/2.5/weather?appid=123MYOWMKEY&lat=-27.4651&lon=153.0231&units=metric
	parameters := url.Values{}
	for key, value := range args {
		parameters.Set(key, value)
	}
	url := BaseURL + api + parameters.Encode()
	return url
}

// Reads a []byte as JSON in the OWNGeoResponse. Tries to find the correct
// location.
func ParseGeoAPIResponse(location string, data []byte) (OWMGeoAPIResponse, error) {
	var responses []OWMGeoAPIResponse
	err := json.Unmarshal(data, &responses)
	if err != nil {
		return OWMGeoAPIResponse{}, fmt.Errorf("invalid API response %s: %w", data, err)
	}
	if len(responses) < 1 {
		return OWMGeoAPIResponse{}, fmt.Errorf("invalid API response %s needs one element", data)
	}

	// Search for the user supplied city,location combination.
	values := strings.Split(location, ",")
	for _, response := range responses {
		if strings.Contains(response.Name, values[0]) && response.Country == values[1] {
			err = ValidateGeoAPIResponse(data, response)
			if err != nil {
				return OWMGeoAPIResponse{}, fmt.Errorf("%w", err)
			}
			return response, nil
		}
	}
	return OWMGeoAPIResponse{}, fmt.Errorf("no match for %s found in %s", location, data)
}

func ValidateGeoAPIResponse(data []byte, response OWMGeoAPIResponse) error {
	if response.Name == "" {
		return fmt.Errorf("invalid API data for name: %s", response.Name)
	}
	if response.Lat < -90 || response.Lat > 90 {
		return fmt.Errorf("invalid API response %s: latitude %v out of range", data, response.Lat)
	}
	if response.Lon < -180 || response.Lon > 180 {
		return fmt.Errorf("invalid API response %s: longitude %v out of range", data, response.Lon)
	}
	if response.Lat == 0 && response.Lon == 0 {
		return fmt.Errorf("invalid API response %s: lat/lon both zero, likely missing data", data)
	}
	return nil
}
