package weather_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"weather"

	owm "github.com/briandowns/openweathermap"
	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

const city = "Brisbane"
const countryCode = "AU"
const BaseURL = "https://api.openweathermap.org"

// Returns the API key associated with environment variable OWM_API_KEY.
func apiKey(t *testing.T) string {
	t.Helper()

	key := os.Getenv("OWM_API_KEY")
	if key == "" {
		fmt.Println("Please set the environment variable OWM_API_KEY.")
		return ""
	}
	return key
}

// Returns the query parameters.
func queryParameters(t *testing.T) string {
	t.Helper()

	// Create proper URL query parameters.
	parameters := url.Values{}
	parameters.Set("q", city+","+countryCode)
	parameters.Set("appid", apiKey(t))

	return parameters.Encode()
}

// Returns the long/lat coordinates of Brisbane.
func coordinates(t *testing.T) owm.Coordinates {
	t.Helper()

	return owm.Coordinates{
		Longitude: 153.0281,
		Latitude:  -27.4679,
	}
}

func Test_NewClientCreatesClientStruct(t *testing.T) {
	t.Parallel()

	client := weather.NewClient(apiKey(t))
	got := client.APIKey()

	want := apiKey(t)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientSetAndGetCityWorks(t *testing.T) {
	t.Parallel()

	city := "Brisbane"
	client := weather.NewClient(apiKey(t))
	client.SetCity(city)
	got := client.City()

	want := city

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientSetAndGetCountryCodeWorks(t *testing.T) {
	t.Parallel()

	countryCode := "AU"
	client := weather.NewClient(apiKey(t))
	client.SetCountryCode(countryCode)
	got := client.CountryCode()

	want := countryCode

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientSetAndGetAPIKeyWorks(t *testing.T) {
	t.Parallel()

	key := "123THISISMYAPIKEY"
	client := weather.NewClient(apiKey(t))
	client.SetAPIKey(key)
	got := client.APIKey()

	want := key

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientFormatWeatherURLWorks(t *testing.T) {
	t.Parallel()

	key := apiKey(t)
	client := weather.NewClient(key)

	got := fmt.Sprintf(
		"%s/data/2.5/weather?appid=%s&q=%s",
		BaseURL,
		client.APIKey(),
		// Avoids hardcoding %2C for URL encoded comma
		url.QueryEscape(city+","+countryCode),
	)

	parameters := queryParameters(t)
	want := fmt.Sprintf(
		"%s/data/2.5/weather?%s",
		BaseURL,
		parameters,
	)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientMakeAPIRequestCreatesWeatherEndpointStruct(t *testing.T) {
	t.Parallel()

	endpoint := "weather"
	client := weather.NewClient(apiKey(t))
	// Need to manually set as we aren't running with command line which injects
	// defaults.
	client.SetCity(city)
	client.SetCountryCode(countryCode)

	var weather weather.WeatherEndpoint
	err := client.MakeAPIRequest(endpoint, &weather)
	if err != nil {
		t.Errorf("error calling MakeAPIRequest: %s", err)
	}
	got := weather.Coordinates

	want := coordinates(t)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// Taken from the-power-of-go-tools/chapter04/count-pflag/count_test.go - see for comments
func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "test/scripts",
		Setup: func(env *testscript.Env) error {
			// Pass the enivonrment variable from host through to test runner.
			key := os.Getenv("OWM_API_KEY")
			env.Setenv("OWM_API_KEY", key)

			// Shared test values
			env.Setenv(
				"BRISBANE_COORDS",
				`coord":{"lon":153.0281,"lat":-27.4679}`,
			)

			env.Setenv(
				"OWM_404_ERROR",
				`unexpected status 404`,
			)

			env.Setenv(
				"OWM_PFLAG_ERROR",
				`error: -c|--city and -C|--country-code must be passed together`,
			)

			return nil
		},
	})
}

// Taken from the-power-of-go-tools/chapter04/count-pflag/count_test.go - see for comments
func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"owm": weather.Main,
	})
}
