package weather_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"weather"

	owm "github.com/briandowns/openweathermap"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rogpeppe/go-internal/testscript"
)

var BrisbaneLatLong = weather.LatLong{Lat: -27.4698, Long: 153.0251}
var BrisbaneLatLongOWM = owm.Coordinates{
	Longitude: 153.0251,
	Latitude:  -27.4698,
}

const BaseURL string = "https://api.openweathermap.org"
const LatLongEndpoint string = "/data/2.5/weather"
const CityCountryEndpoint string = "/geo/1.0/direct"
const ZipCountryEndpoint string = "/geo/1.0/zip"

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

// Returns a string with the URL to make an API call.
func buildURL(t *testing.T, client weather.Client) string {
	t.Helper()

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

func Test_ClientBuildURLWorks(t *testing.T) {
	t.Parallel()

	client := weather.Client{
		Key:      apiKey(t),
		Location: BrisbaneLatLong,
	}
	got := client.BuildURL()
	want := buildURL(t, client)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_LatLongEndpointReturnsCorrectly(t *testing.T) {
	t.Parallel()

	want := LatLongEndpoint
	got := weather.LatLong{}.Endpoint()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_NewClientCreatesClientStruct(t *testing.T) {
	t.Parallel()

	client := weather.NewClient(apiKey(t))
	got := client.Key

	want := apiKey(t)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ClientMakeAPIRequestCreatesWeatherEndpointStruct(t *testing.T) {
	t.Parallel()

	client := weather.NewClient(apiKey(t))
	// Need to manually set as we aren't running with command line which injects
	// defaults.
	client.Location = weather.LatLong{Lat: -27.4679, Long: 153.0281}

	url := buildURL(t, *client)
	var weather weather.WeatherEndpoint
	err := client.MakeAPIRequest(url, &weather)
	if err != nil {
		t.Errorf("error calling MakeAPIRequest: %s", err)
	}
	got := weather.Coordinates

	want := BrisbaneLatLongOWM

	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.01)) {
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
				"OWM_OUTPUT",
				`"name":"Brisbane"`,
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
