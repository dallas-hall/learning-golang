package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"weather"
)

const Usage = `Usage: weather LOCATION

eg: ` + weather.BrisbaneLocation

func main() {
	// Parse command line
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		return
	}
	location := os.Args[1]

	// Get API key
	key, err := weather.APIKey()
	if err != nil {
		log.Fatal(err)
	}

	// Create query parameters to lookup lat/lon by city,country
	// https://openweathermap.org/api/geocoding-api?collection=other#direct
	args := map[string]string{
		"q":     location,
		"appid": key,
	}

	// Make API call to lookup lat/lon by city,country
	URL := weather.FormatURL(weather.GeoAPI, args)
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected HTTP response status: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	coordinates, err := weather.ParseGeoAPIResponse(location, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(coordinates)

	// Delete our Geo API query params and populate Data API ones.
	delete(args, "q")

	// Create query parameters to lookup weather by lat/lon
	// https://openweathermap.org/api/current?collection=current_forecast#geo
	args["lat"] = fmt.Sprintf("%.4f", coordinates.Lat)
	args["lon"] = fmt.Sprintf("%.4f", coordinates.Lon)
	args["units"] = "metric"

	// Make API call to lookup weather by lat/lon
	URL = weather.FormatURL(weather.DataAPI, args)
	response, err = http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("unexpected HTTP response status: %s", response.Status)
	}

	data, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	conditions, err := weather.ParseDataAPIResponse(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conditions)
}
