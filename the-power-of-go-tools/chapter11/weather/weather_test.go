package weather_test

import (
	"net/url"
	"os"
	"strings"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
)

func Test_APIKeyReturnsCorrectly(t *testing.T) {
	want := "737"
	got, err := weather.APIKey()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(got, want) {
		t.Fatalf("%q not found in %q", want, got)
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

	}

	want := weather.Conditions{
		Summary: "Clear",
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
		t.Fatal(err)
	}

	want := weather.OWMGeoAPIResponse{
		Name:    "City of Brisbane",
		Country: "AU",
		Lat:     weather.BrisbaneLatitude,
		Lon:     weather.BrisbaneLongitude,
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}

}

func Test_ParseDataAPIResponseReturnsErrorFromEmptyData(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseDataAPIResponse([]byte{})
	if err == nil {
		t.Fatal("want error parsing empty response, got nil")
	}
}

func Test_ParseGeoAPIResponseReturnsErrorFromEmptyData(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseGeoAPIResponse(weather.BrisbaneLocation, []byte{})
	if err == nil {
		t.Fatal("want error parsing empty response, got nil")
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
		t.Fatal("want error parsing invalid response, got nil")
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
		t.Fatal("want error parsing invalid response, got nil")
	}
}

func Test_ParseGeoAPIResponseReturnsErrorFromInvalidLocation(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseGeoAPIResponse("", []byte{})
	if err == nil {
		t.Fatal("want error parsing invalid location, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidName(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:    "",
		Country: "AU",
		Lat:     weather.BrisbaneLatitude,
		Lon:     weather.BrisbaneLongitude,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Fatal("want error parsing invalid name, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLat(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:    "Brisbane",
		Country: "AU",
		Lat:     91,
		Lon:     weather.BrisbaneLongitude,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Fatal("want error parsing invalid latitude, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLon(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:    "Brisbane",
		Country: "AU",
		Lat:     weather.BrisbaneLatitude,
		Lon:     -181,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Fatal("want error parsing invalid longitude, got nil")
	}
}

func Test_ValidateGeoAPIResponseReturnsErrorWithInvalidLatAndLon(t *testing.T) {
	t.Parallel()

	data := []byte{}
	response := weather.OWMGeoAPIResponse{
		Name:    "Brisbane",
		Country: "AU",
		Lat:     0,
		Lon:     0,
	}

	err := weather.ValidateGeoAPIResponse(data, response)
	if err == nil {
		t.Fatal("want error parsing invalid latitude and longitude, got nil")
	}
}

func Test_FormatURLReturnsCorrectString(t *testing.T) {
	t.Parallel()

	api := "/data/2.5/weather?"
	lat := "-27.4651"
	long := "153.0231"
	units := "&units=metric"
	key := "123MYOWMKEY"
	want := weather.BaseURL + api + "appid=" + key + "&lat=" + lat +
		"&lon=" + long + units

	args := map[string]string{
		"lat":   lat,
		"lon":   long,
		"appid": key,
		"units": "metric",
	}

	encoded := weather.FormatURL(api, args)
	got, err := url.QueryUnescape(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
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
