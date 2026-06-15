package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var defaultHTTPClient = &http.Client{}

const maxErrorBodyPreview = 256

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

		reqBody := bytes.NewReader(bodyBytes)

		req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
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

			if err := wait(ctx, backoff); err != nil {
				return nil, err
			}

			backoff *= 2
			continue
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			if attempt == defaultMaxRetries {
				return nil, readErr
			}

			log.Warn().
				Err(readErr).
				Int("attempt", attempt+1).
				Msg("failed reading response body, retrying")

			if err := wait(ctx, backoff); err != nil {
				return nil, err
			}

			backoff *= 2
			continue
		}

		bodyText := strings.TrimSpace(string(respBody))

		// ---------- 429 RATE LIMIT ----------
		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt == defaultMaxRetries {
				return nil, fmt.Errorf("http %d: %s", resp.StatusCode, truncate(bodyText))
			}

			retryDelay := parseRetryAfter(resp.Header.Get("Retry-After"), backoff)

			log.Warn().
				Int("status", resp.StatusCode).
				Int("attempt", attempt+1).
				Dur("retry_after", retryDelay).
				Msg("rate limited, retrying")

			if err := wait(ctx, retryDelay); err != nil {
				return nil, err
			}

			backoff *= 2
			continue
		}

		// ---------- 5xx RETRIES ----------
		if resp.StatusCode >= 500 {
			if attempt == defaultMaxRetries {
				return nil, fmt.Errorf("http %d: %s", resp.StatusCode, truncate(bodyText))
			}

			log.Warn().
				Int("status", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("server error, retrying")

			if err := wait(ctx, backoff); err != nil {
				return nil, err
			}

			backoff *= 2
			continue
		}

		// ---------- NON-2XX FAIL FAST ----------
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("http %d: %s", resp.StatusCode, truncate(bodyText))
		}

		// Success: restore body for caller
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
		return resp, nil
	}

	return nil, fmt.Errorf("exceeded max retries")
}

func wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	}
}

func parseRetryAfter(h string, fallback time.Duration) time.Duration {
	if h == "" {
		return fallback
	}

	// delta-seconds
	if seconds, err := strconv.Atoi(strings.TrimSpace(h)); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// HTTP-date
	if t, err := http.ParseTime(h); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}

	return fallback
}

func truncate(s string) string {
	if len(s) <= maxErrorBodyPreview {
		return s
	}
	return s[:maxErrorBodyPreview] + "...(truncated)"
}
