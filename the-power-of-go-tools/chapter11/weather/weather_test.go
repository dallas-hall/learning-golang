package weather_test

import (
	"net/url"
	"os"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseDataAPIResponseReturnsGoStruct(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("test/data/sample.json")
	if err != nil {
		t.Fatal(err)
	}
	want := weather.Conditions{
		Summary: "Clear",
	}

	got, err := weather.ParseDataAPIResponse(data)
	if err != nil {

	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func Test_ParseDataAPIResponseReturnsErrorFromEmptyData(t *testing.T) {
	t.Parallel()

	_, err := weather.ParseDataAPIResponse([]byte{})
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
