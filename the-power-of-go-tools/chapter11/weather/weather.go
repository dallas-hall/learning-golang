package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type OWMDataAPIResponse struct {
	Name  string `json:"name"`
	Coord struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	} `json:"coord"`
	Main struct {
		Temperature float64 `json:"temp"`
		FeelsLike   float64 `json:"feels_like"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
}

type OWMGeoAPIResponse struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type Conditions struct {
	Summary     string
	Temperature Temperature
	FeelsLike   Temperature
}

type OWMClient struct {
	Key        string
	BaseURL    string
	DataAPI    string
	GeoAPI     string
	HTTPClient *http.Client
}

const BrisbaneLocation = "Brisbane,AU"
const BrisbaneLatitude float64 = -27.4689682
const BrisbaneLongitude float64 = 153.0234991
const Usage = `Usage: weather LOCATION

eg: ` + BrisbaneLocation

// Create a new type so we can create a method for it.
type Temperature float64

func (t Temperature) Celsius() float64 {
	return float64(t) - 273.15
}

func NewClient(key string) *OWMClient {
	c := &OWMClient{
		Key:     key,
		BaseURL: "https://api.openweathermap.org",
		DataAPI: "/data/2.5/weather",
		GeoAPI:  "/geo/1.0/direct",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return c
}

// Returns the API key associated with environment variable OWM_API_KEY.
func APIKey() (string, error) {
	key := os.Getenv("OWM_API_KEY")
	if key == "" {
		msg := "environment variable OWM_API_KEY is not set"
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
		return Conditions{}, fmt.Errorf("invalid API response %s with %w", data, err)
	}
	if len(response.Weather) < 1 {
		return Conditions{}, fmt.Errorf("invalid API response %s needs one Weather element", data)
	}
	conditions := Conditions{
		Summary: response.Weather[0].Main,
		// We need to cast this because of our custom type used for Celsius()
		Temperature: Temperature(response.Main.Temperature),
		FeelsLike:   Temperature(response.Main.FeelsLike),
	}
	return conditions, nil
}

// Returns a URL using url.Values.Encode for encoding query parameters.
// Requires API endpoint to call, and the args map containing the key.
func (c OWMClient) FormatURL(api string, args map[string]string) string {
	// https://api.openweathermap.org/data/2.5/weather?appid=123MYOWMKEY&lat=-27.4651&lon=153.0231&units=metric
	parameters := url.Values{}
	for key, value := range args {
		parameters.Set(key, value)
	}
	return c.BaseURL + api + "?" + parameters.Encode()
}

// Reads a []byte as JSON in the OWNGeoResponse. Tries to find the correct
// location.
func ParseGeoAPIResponse(location string, data []byte) (OWMGeoAPIResponse, error) {
	if !strings.Contains(location, ",") {
		return OWMGeoAPIResponse{}, fmt.Errorf("invalid location %q must be %q", location, "City,CountryCode")
	}

	var responses []OWMGeoAPIResponse
	err := json.Unmarshal(data, &responses)
	if err != nil {
		return OWMGeoAPIResponse{}, fmt.Errorf("invalid API response %s with %w", data, err)
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
		return fmt.Errorf("invalid API data for name: %q", response.Name)
	}
	if response.Latitude < -90 || response.Latitude > 90 {
		return fmt.Errorf("invalid API response %s: latitude %v out of range", data, response.Latitude)
	}
	if response.Longitude < -180 || response.Longitude > 180 {
		return fmt.Errorf("invalid API response %s: longitude %v out of range", data, response.Longitude)
	}
	if response.Latitude == 0 && response.Longitude == 0 {
		return fmt.Errorf("invalid API response %s: lat/lon both zero, likely missing data", data)
	}
	return nil
}

func (c OWMClient) GetWeather(location string) (Conditions, error) {
	// Create query parameters to lookup lat/lon by city,country
	// https://openweathermap.org/api/geocoding-api?collection=other#direct
	args := map[string]string{
		"q":     location,
		"appid": c.Key,
	}

	// Make API call to lookup lat/lon by city,country
	URL := c.FormatURL(c.GeoAPI, args)
	response, err := c.HTTPClient.Get(URL)
	if err != nil {
		return Conditions{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf(
			"unexpected HTTP response status: %s",
			response.Status,
		)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Conditions{}, err
	}

	coordinates, err := ParseGeoAPIResponse(location, data)
	if err != nil {
		return Conditions{}, err
	}

	// Delete our Geo API query params and populate Data API ones.
	delete(args, "q")

	// Create query parameters to lookup weather by lat/lon
	// https://openweathermap.org/api/current?collection=current_forecast#geo
	args["lat"] = fmt.Sprintf("%.4f", coordinates.Latitude)
	args["lon"] = fmt.Sprintf("%.4f", coordinates.Longitude)
	//args["units"] = "metric"

	// Make API call to lookup weather by lat/lon
	URL = c.FormatURL(c.DataAPI, args)
	response, err = c.HTTPClient.Get(URL)
	if err != nil {
		return Conditions{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf(
			"unexpected HTTP response status: %s",
			response.Status,
		)
	}

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return Conditions{}, err
	}

	conditions, err := ParseDataAPIResponse(data)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, nil
}

// Conveience wrapper to create a client with the supplied key and then
// return the weather summary for the given CITY,COUNTRYCODE location.
func Get(location, key string) (Conditions, error) {
	c := NewClient(key)
	conditions, err := c.GetWeather(location)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, err
}

func Main() {
	// Parse command line
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		return
	}
	location := os.Args[1]

	// Get API key
	key, err := APIKey()
	if err != nil {
		log.Fatal(err)
	}

	// Create the client
	client := NewClient(key)

	conditions, err := client.GetWeather(location)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"Summary: %s\nTemperature: %.1f\nFeels like: %.1f\n",
		conditions.Summary,
		conditions.Temperature.Celsius(),
		conditions.FeelsLike.Celsius(),
	)
}
