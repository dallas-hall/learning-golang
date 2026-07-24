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
	"strconv"
	"strings"

	// https://pkg.go.dev/github.com/spf13/pflag
	flag "github.com/spf13/pflag"

	owm "github.com/briandowns/openweathermap"
)

// The URL has the full path to the API endpoint. Must be supplied to the
// constructor NewClient().
// The API Key can either be added to the URI or elsewhere (e.g. HTTP header)
// The location holds an interface of accepted API query formats.
type Client struct {
	Key      string
	Location LocationQuery
}

// The client needs a general interface for all location API calls it can make.
// We did this so we could handle Lat/Long, City/Country, Zip/Country. But the
// last 2 were removed due to complexity and I couldn't be fucked spending more
// time on this exercise which turned into a mini project. LOL.
type LocationQuery interface {
	QueryValues() map[string]string
	Endpoint() string
}

type LatLong struct {
	Lat, Long float64
}

func ParseLatLong(s string) (LatLong, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return LatLong{}, fmt.Errorf("expected lat,long format, got %q", s)
	}

	// We are wrapping any strconv.Parsefloat error with %w
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return LatLong{}, fmt.Errorf("invalid latitude %q: %w", parts[0], err)
	}

	long, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return LatLong{}, fmt.Errorf("invalid longitude %q: %w", parts[1], err)
	}

	return NewLatLong(lat, long)
}

func NewLatLong(lat, long float64) (LatLong, error) {
	if lat < -90 || lat > 90 {
		return LatLong{}, fmt.Errorf("latitude %v out of range [-90, 90]", lat)
	}
	if long < -180 || long > 180 {
		return LatLong{}, fmt.Errorf("longitude %v out of range [-180, 180]", long)
	}
	return LatLong{Lat: lat, Long: long}, nil
}

func (l LatLong) QueryValues() map[string]string {
	return map[string]string{
		"lat": fmt.Sprintf("%.4f", l.Lat),
		"lon": fmt.Sprintf("%.4f", l.Long),
	}
}

func (l LatLong) Endpoint() string {
	return "/data/2.5/weather"
}

// Create a new client with the supplied URL.
func NewClient(key string) *Client {
	return &Client{
		Key: key,
	}
}

func (client Client) BuildURL() string {
	// Create proper URL query parameters.
	parameters := url.Values{}
	for key, value := range client.Location.QueryValues() {
		parameters.Set(key, value)
	}
	parameters.Set("appid", client.Key)

	return fmt.Sprintf(
		"%s%s?%s",
		BaseURL,
		client.Location.Endpoint(),
		parameters.Encode(),
	)
}

// Pass in a complete URL to the API endpoint and something to hold the JSON
// data returned by the API request. An error is returned if JSON data cannot
// be returned.
func (client *Client) MakeAPIRequest(url string, result any) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// result is passed in pointer
	err = json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("%v in %q", err, data)
	}

	return nil
}

// Holds the OpenWeatherMap data from the "weather" API end point.
// This uses "github.com/briandowns/openweathermap" which at the time of writing
// didn't map all the fields being returned by the API. Those fields are lost.
// The API also only returns data relevant to the current weather, e.g. if it
// isn't raining then there will be no "rain" returned.
// https://openweathermap.org/api/current?collection=current_forecast#fields_json
type WeatherEndpoint struct {
	Coordinates owm.Coordinates `yaml:"coord" json:"coord"`
	Weather     []owm.Weather   `yaml:"weather" json:"weather"`
	Base        string          `yaml:"base" json:"base"`
	Main        owm.Main        `yaml:"main" json:"main"`
	Visibility  int             `yaml:"visibility" json:"visibility"`
	Wind        owm.Wind        `yaml:"wind" json:"wind"`
	Clouds      owm.Clouds      `yaml:"clouds" json:"clouds"`
	Rain        owm.Rain        `yaml:"rain" json:"rain"`
	Snow        owm.Snow        `yaml:"snow" json:"snow"`
	DT          int             `yaml:"dt" json:"dt"`
	Sys         owm.Sys         `yaml:"sys" json:"sys"`
	Timezone    int             `yaml:"timezone" json:"timezone"`
	ID          int             `yaml:"id" json:"id"`
	Name        string          `yaml:"name" json:"name"`
	Cod         int             `yaml:"cod" json:"cod"`
}

const BaseURL = "https://api.openweathermap.org"

// Returns the API key associated with environment variable OWM_API_KEY.
func APIKey() (string, error) {
	key := os.Getenv("OWM_API_KEY")
	if key == "" {
		msg := "error: environment variable OWM_API_KEY is not set"
		return "", errors.New(msg)
	}
	return key, nil
}

func Main() {
	// Taken from the-power-of-go-tools/chapter04/count-pflag/count.go see there for comments.
	latLongFlag := flag.StringP(
		"lat-long",
		"l",
		"-27.4698,153.0251",
		"Search using latitude and longitude delimited by comma, eg -27.4698,153.0251",
	)

	// Update the -h output.
	flag.Usage = func() {
		fmt.Printf("Usage: %s option\n", os.Args[0])
		fmt.Println("Create a client and send API requests to OpenWeatherMap.")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	// Check the command line for arguments and assign them to our matching flags.
	// This stops parsing as soon as it see a non-flag arg.
	flag.Parse()

	// Configure client
	key, err := APIKey()
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(key)

	// Make API call
	client.Location, err = ParseLatLong(*latLongFlag)
	if err != nil {
		log.Fatal(err)
	}

	var weather WeatherEndpoint
	url := client.BuildURL()

	err = client.MakeAPIRequest(url, &weather)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(weather)
	if err != nil {
		log.Fatalf("error: unable to json.Marshal %v because %s\n", weather, err)
	}
	fmt.Println(string(data))
}
