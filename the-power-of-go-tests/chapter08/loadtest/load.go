package loadtest

import (
	"io"
	"net/http"
	"os"
)

type LoadTester struct {
	URL                 string
	Output, ErrorOutput io.Writer
	HTTPClient          *http.Client
	Fanout              int
	UserAgent           string
}

// There are too many parameters being passed in here. There should be sane
// defaults that users can choose to overwrite if they want to.
func NewLoadTesterBad(URL string, stdout, stderr io.Writer, httpClient *http.Client, fanout int, userAgent string) *LoadTester {
	return &LoadTester{
		URL:         URL,
		Output:      stdout,
		ErrorOutput: stderr,
		HTTPClient:  httpClient,
		Fanout:      fanout,
		UserAgent:   userAgent,
	}
}

// Use functional options to override sane defaults, rather than passing in
// every parameter. See `the-power-of-go-tools/chapter03/count/count.go`for
// details.
func NewLoadTester(URL string, opts ...TesterOption) *LoadTester {
	lt := &LoadTester{
		URL:         URL,
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
		HTTPClient:  http.DefaultClient,
		Fanout:      10,
		UserAgent:   "load v2 (with functional options)",
	}
	for _, opt := range opts {
		opt(lt)
	}
	return lt
}

type TesterOption func(*LoadTester)

func WithOutput(w io.Writer) TesterOption {
	return func(lt *LoadTester) {
		lt.Output = w
	}
}

func WithErrorOutput(w io.Writer) TesterOption {
	return func(lt *LoadTester) {
		lt.ErrorOutput = w
	}
}

func WithHTTPClient(c *http.Client) TesterOption {
	return func(lt *LoadTester) {
		lt.HTTPClient = c
	}
}

func WithConcurrentRequests(fanout int) TesterOption {
	return func(lt *LoadTester) {
		lt.Fanout = fanout
	}
}

func WithUserAgent(agent string) TesterOption {
	return func(lt *LoadTester) {
		lt.UserAgent = agent
	}
}
