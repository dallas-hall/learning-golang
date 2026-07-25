package weather_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const fakeKey = "123-MY-FAKE-KEY"

func Test_APIKeyReturnsCorrectly(t *testing.T) {
	want := "737"
	got, err := weather.APIKey()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got, want) {
		t.Errorf("%q not found in %q", want, got)
	}
}

func Test_ParseDataAPIResponseReturnsGoStruct(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/data-api-sample.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := weather.ParseDataAPIResponse(data)
	if err != nil {
		t.Fatal(err)
	}

	want := weather.Conditions{
		Summary:     "Clear",
		Temperature: 293.79,
		FeelsLike:   293.69,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ParseGeoAPIResponseReturnsGoStruct(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/geo-api-sample.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := weather.ParseGeoAPIResponse(weather.BrisbaneLocation, data)
	if err != nil {
		t.Error(err)
	}

	want := weather.OWMGeoAPIResponse{
		Name:      "City of Brisbane",
		Country:   "AU",
		Latitude:  weather.BrisbaneLatitude,
		Longitude: weather.BrisbaneLongitude,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func Test_ParseDataAPIResponseReturnsErrorFromEmptyData(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseDataAPIResponse([]byte{})
	if err == nil {
		t.Error("want error parsing empty response, got nil")
	}
}

func Test_ParseGeoAPIResponseReturnsErrorFromEmptyData(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseGeoAPIResponse(weather.BrisbaneLocation, []byte{})
	if err == nil {
		t.Error("want error parsing empty response, got nil")
	}
}

func Test_ParseDataAPIResponseReturnsErrorFromInvalidJSON(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/sample.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = weather.ParseDataAPIResponse(data)
	if err == nil {
		t.Error("want error parsing invalid response, got nil")
	}
}

func Test_ParseGeoAPIResponseReturnsErrorFromInvalidJSON(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/sample.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = weather.ParseGeoAPIResponse(weather.BrisbaneLocation, data)
	if err == nil {
		t.Error("want error parsing invalid response, got nil")
	}
}

func Test_ParseGeoAPIResponseReturnsErrorFromInvalidLocation(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseGeoAPIResponse("", []byte{})
	if err == nil {
		t.Error("want error parsing invalid location, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidName(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:      "",
		Country:   "AU",
		Latitude:  weather.BrisbaneLatitude,
		Longitude: weather.BrisbaneLongitude,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Error("want error parsing invalid name, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLat(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:      "Brisbane",
		Country:   "AU",
		Latitude:  91,
		Longitude: weather.BrisbaneLongitude,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Error("want error parsing invalid latitude, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLon(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:      "Brisbane",
		Country:   "AU",
		Latitude:  weather.BrisbaneLatitude,
		Longitude: -181,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Error("want error parsing invalid longitude, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLatAndLon(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:      "Brisbane",
		Country:   "AU",
		Latitude:  0,
		Longitude: 0,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Error("want error parsing invalid latitude and longitude, got nil")
	}
}

func Test_FormatURLReturnsCorrectString(t *testing.T) {
	t.Parallel()

	baseURL := "https://api.openweathermap.org"
	api := "/data/2.5/weather?"
	lat := "-27.4651"
	long := "153.0231"
	units := "&units=metric"
	want := baseURL + api + "appid=" + fakeKey + "&lat=" + lat +
		"&lon=" + long + units

	client := weather.NewClient(fakeKey)

	args := map[string]string{
		"lat":   lat,
		"lon":   long,
		"appid": fakeKey,
		"units": "metric",
	}

	encoded := client.FormatURL(client.DataAPI, args)
	got, err := url.QueryUnescape(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_HTTPGetFromLocalServerWorks(t *testing.T) {
	t.Parallel()

	// Create a test server with a self signed TLS certificate
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "test/data/geo-api-sample.json")
			}))
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	want := http.StatusOK
	got := response.StatusCode

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_HTTPSGetFromLocalServerWorks(t *testing.T) {
	t.Parallel()

	// Create a test server with a self signed TLS certificate
	testServer := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "test/data/data-api-sample.json")
			}))
	defer testServer.Close()

	// This allows us to test using TLS with the self signed certificate
	client := testServer.Client()
	response, err := client.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	want := http.StatusOK
	got := response.StatusCode

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_GetWeatherWorksWithLocalHTTPSServer(t *testing.T) {
	t.Parallel()

	wc := weather.NewClient(fakeKey)

	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Recreate the OpenWeather Map Geo and Data API endpoints.
				switch {
				case strings.HasPrefix(r.URL.Path, wc.GeoAPI):
					http.ServeFile(w, r, "test/data/geo-api-sample.json")
				case strings.HasPrefix(r.URL.Path, wc.DataAPI):
					http.ServeFile(w, r, "test/data/data-api-sample.json")
				default:
					http.NotFound(w, r)
				}
			}))
	defer ts.Close()

	wc.BaseURL = ts.URL
	wc.HTTPClient = ts.Client()

	got, err := wc.GetWeather(weather.BrisbaneLocation)
	if err != nil {
		t.Fatal(err)
	}

	want := weather.Conditions{
		Summary:     "Clear",
		Temperature: 293.79,
		FeelsLike:   293.69,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_KelvinToCelsiusConversionWorks(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/data-api-sample.json")
	if err != nil {
		t.Fatal(err)
	}
	conditions, err := weather.ParseDataAPIResponse(data)
	if err != nil {
		t.Fatal(err)
	}

	want := 20.64
	got := conditions.Temperature.Celsius()

	// We need a tolerance as "20.640000000000043" can be returned instead of 20.64
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.0001)) {
		t.Error(cmp.Diff(want, got))
	}
}

// func Test(t *testing.T) {
// 	t.Parallel()
// 	testscript.Run(t, testscript.Params{
// 		Dir: "test/scripts",
// 		Setup: func(env *testscript.Env) error {
// 			// Pass the enivonrment variable from host through to test runner.
// 			key := os.Getenv("OWM_API_KEY")
// 			env.Setenv("OWM_API_KEY", key)
// 			return nil
// 		},
// 	})
// }
