package request

import (
	"errors"
	"net/http"
)

// Create a sentinel error, a named error value. Just remember you will lose
// error information if you don't wrap the original error within the sentinel
// error. e.g. our sentinel error doesn't have the HTTP status code.
var ErrRateLimit = errors.New("rate limit")

func Request(URL string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusTooManyRequests {
		return ErrRateLimit
	}
	return nil
}
