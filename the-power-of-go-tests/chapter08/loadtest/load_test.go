package loadtest_test

import (
	"bytes"
	"loadtest"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewLoadTesterBadSetsConfigurationCorrectly(t *testing.T) {
	t.Parallel()
	want := &loadtest.LoadTester{
		URL:         "https://example.com",
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
		HTTPClient:  http.DefaultClient,
		Fanout:      10,
		UserAgent:   "loadtest",
	}
	got := loadtest.NewLoadTesterBad(
		"https://example.com",
		os.Stdout,
		os.Stdin,
		http.DefaultClient,
		10,
		"loadtest",
	)
	// There are unexported fields inside of os.File file pointer (eg name) and
	// they are outside of this package. So we tell cmp to ignore those errors.
	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(os.File{})) {
		t.Error(cmp.Diff(want, got, cmpopts.IgnoreUnexported(os.File{})))
	}
}

func TestNewLoadTesterSetsConfigurationCorrectly(t *testing.T) {
	t.Parallel()
	buffer := new(bytes.Buffer)
	want := &loadtest.LoadTester{
		URL:         "https://example.com",
		Output:      buffer,
		ErrorOutput: buffer,
		HTTPClient:  &http.Client{Timeout: 3 * time.Second},
		Fanout:      10,
		UserAgent:   "loadtest",
	}
	got := loadtest.NewLoadTester("https://example.com",
		loadtest.WithOutput(buffer),
		loadtest.WithErrorOutput(buffer),
		loadtest.WithHTTPClient(&http.Client{Timeout: 3 * time.Second}),
		// Skip fanout to show the default is used
		loadtest.WithUserAgent("loadtest"),
	)
	ignore := cmpopts.IgnoreUnexported(os.File{}, bytes.Buffer{})
	if !cmp.Equal(want, got, ignore) {
		t.Error(cmp.Diff(want, got, ignore))
	}
}
