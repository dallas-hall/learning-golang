package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type OWMResponse struct {
	Weather []struct {
		Main        string
		Description string
	}
}

type Conditions struct {
	Summary string
}

const BaseURL = "https://api.openweathermap.org"
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

func ParseResponse(data []byte) (Conditions, error) {
	var response OWMResponse
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

func BasicPrint() {
	key, err := APIKey()
	if err != nil {
		log.Fatal(err)
	}

	api := "/data/2.5/weather?"
	lat := BrisbaneLatitude
	long := BrisbaneLongitude
	units := "&units=metric"
	key = "&appid=" + key
	url := BaseURL + api + "lat=" + lat + "&lon=" + long + key + units
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status %d", response.StatusCode)
	}

	io.Copy(os.Stdout, response.Body)
}

func UpdatedPrint() {
	key, err := APIKey()
	if err != nil {
		log.Fatal(err)
	}

	api := "/data/2.5/weather?"
	args := map[string]string{
		"lat":   BrisbaneLatitude,
		"lon":   BrisbaneLongitude,
		"appid": key,
		"units": "metric",
	}

	parameters := url.Values{}
	for key, value := range args {
		parameters.Set(key, value)
	}
	url := BaseURL + api + parameters.Encode()
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	conditions, err := ParseResponse(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(conditions)
}

func FormatURL(api string, args map[string]string, key string) string {
	return ""
}

func main() {
	BasicPrint()
	fmt.Println()
	UpdatedPrint()
}
