package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var defaultHTTPClient = &http.Client{}

func DoHTTPRequest(ctx context.Context, method string, data any, url string) (*http.Response, error) {
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

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return resp, nil
}

const (
	defaultMaxRetries = 5
	initialBackoff    = 200 * time.Millisecond
)

func DoRetryableHTTPRequest(ctx context.Context, method string, data any, url string) (*http.Response, error) {
	var bodyBytes []byte
	var err error

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		bodyBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	backoff := initialBackoff

	for attempt := 0; attempt <= defaultMaxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		var body *bytes.Reader
		if bodyBytes != nil {
			body = bytes.NewReader(bodyBytes)
		} else {
			body = bytes.NewReader(nil)
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		log.Info().
			Int("attempt", attempt+1).
			Str("method", method).
			Str("url", url).
			Msg("making HTTP request")

		resp, err := defaultHTTPClient.Do(req)

		if err != nil {
			if attempt == defaultMaxRetries {
				return nil, err
			}

			log.Warn().
				Err(err).
				Int("attempt", attempt+1).
				Msg("request failed, retrying")

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			backoff *= 2
			continue
		}

		if resp.StatusCode >= 500 {
			if attempt == defaultMaxRetries {
				return resp, fmt.Errorf("HTTP error: %s", resp.Status)
			}

			resp.Body.Close()

			log.Warn().
				Int("status", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("server error, retrying")

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			backoff *= 2
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return resp, fmt.Errorf("HTTP error: %s", resp.Status)
		}

		return resp, nil
	}
	return nil, fmt.Errorf("exceeded max retries")
}
