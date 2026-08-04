package request_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"request"
	"testing"
)

// Simulate a HTTP server that is rate limiting requests.
func newRateLimitingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
}

func Test_RequestReturnsErrRateLimitErrorWhenRateLimited(t *testing.T) {
	t.Parallel()
	ts := newRateLimitingServer()
	defer ts.Close()
	err := request.Request(ts.URL)
	if !errors.Is(err, request.ErrRateLimit) {
		t.Errorf("wrong error: %v", err)
	}
}
