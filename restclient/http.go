package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// DoHTTPRequest sends an HTTP request with the specified method, data, and URL using the provided context.
// For POST requests, it marshals the data as JSON and includes it in the request body.
// Sets appropriate headers for JSON requests and logs the request details.
// Returns the HTTP response or an error if the request fails or the status is not OK.
func DoHTTPRequest(ctx context.Context, method string, data any, url string) (*http.Response, error) {
	client := &http.Client{}

	var requestBody []byte
	if method == http.MethodPost {
		body, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		requestBody = body
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	log.Info().Msgf("making %s request to: %s", method, req.URL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return resp, nil
}
