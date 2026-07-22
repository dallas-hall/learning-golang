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

	// https://pkg.go.dev/github.com/spf13/pflag
	flag "github.com/spf13/pflag"

	owm "github.com/briandowns/openweathermap"
)

// The URL has the full path to the API endpoint. Must be supplied to the
// constructor NewClient().
// The API Key can either be added to the URI or elsewhere (e.g. HTTP header)
// The city and country code can be used for location based endpoints.
type Client struct {
	apiKey      string
	city        string
	countryCode string
}

// Create a new client with the supplied URL.
func NewClient(key string) *Client {
	return &Client{
		apiKey: key,
	}
}

// Return the city in a new string.
func (client *Client) City() string {
	return client.city
}

// Return the country code  in a new string.
func (client *Client) CountryCode() string {
	return client.countryCode
}

// Return the API key in a new string.
func (client *Client) APIKey() string {
	return client.apiKey
}

// Update the city string.
func (client *Client) SetCity(city string) {
	client.city = city
}

// Update the country code string.
func (client *Client) SetCountryCode(countryCode string) {
	client.countryCode = countryCode
}

// Update the API key string.
func (client *Client) SetAPIKey(key string) {
	client.apiKey = key
}

func (client *Client) FormatWeatherURL() string {
	// Create proper URL query parameters.
	parameters := url.Values{}
	parameters.Set("q", client.city+","+client.countryCode)
	parameters.Set("appid", client.apiKey)

	return fmt.Sprintf(
		"%s/data/2.5/weather?%s",
		BaseURL,
		parameters.Encode(),
	)
}

// Pass in something to hold the JSON data returned by the API request.
// An error is returned if JSON data cannot be returned.
func (client *Client) MakeAPIRequest(endpoint string, result any) error {
	var url string
	switch {
	// TODO: Add more endpoints.
	case strings.ToLower(endpoint) == "weather":
		url = client.FormatWeatherURL()
	default:
		url = client.FormatWeatherURL()
	}

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
type WeatherEndpoint struct {
	Coordinates owm.Coordinates `yaml:"coord" json:"coord"`
	Weather     []owm.Weather   `yaml:"weather" json:"weather"`
	Base        string          `yaml:"base" json:"base"`
	Main        owm.Main        `yaml:"main" json:"main"`
	Visibility  int             `yaml:"visibility" json:"visibility"`
	Wind        owm.Wind        `yaml:"wind" json:"wind"`
	Clouds      owm.Clouds      `yaml:"clouds" json:"clouds"`
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

// Returns the JSON direct from the API.
func (client *Client) PrintWeather() error {
	if client.city == "" || client.countryCode == "" {
		msg := fmt.Sprintf(
			"missing input for city or country code. city:%q countryCode:%q",
			client.city,
			client.countryCode,
		)
		return errors.New(msg)
	}

	url := client.FormatWeatherURL()

	// Send a HTTP client request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the server's response
	if response.StatusCode != http.StatusOK {
		msg := fmt.Sprintf(
			"unexpected HTTP response status %d\n",
			response.StatusCode,
		)
		return errors.New(msg)
	}

	io.Copy(os.Stdout, response.Body)
	return nil
}

func Main() {
	// Taken from the-power-of-go-tools/chapter04/count-pflag/count.go see there for comments.
	cityFlag := flag.StringP("city", "c", "Brisbane", "Set the city")
	countryCodeFlag := flag.StringP("country-code", "C", "AU", "Set the country code.")
	weatherFlag := flag.BoolP("weather-api", "w", true, "Call the Weather API endpoint.")

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

	// Check for valid flag combinations
	if !flag.Lookup("city").Changed && flag.Lookup("country-code").Changed ||
		flag.Lookup("city").Changed && !flag.Lookup("country-code").Changed {
		log.Fatal("error: -c|--city and -C|--country-code must be passed together")
	}

	// Configure client
	key, err := APIKey()
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient(key)
	client.SetCity(*cityFlag)
	client.SetCountryCode(*countryCodeFlag)

	// Make API calls
	switch {
	// TODO: Add more endpoints.
	case *weatherFlag:
		//client.PrintWeather()
		var weather WeatherEndpoint
		err := client.MakeAPIRequest("weather", &weather)
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.Marshal(weather)
		if err != nil {
			log.Fatalf("error: unable to json.Marshal %v because %s\n", weather, err)
		}
		fmt.Println(string(data))
	}
}
