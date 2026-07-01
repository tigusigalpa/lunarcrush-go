package lunarcrush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// queryValues builds a url.Values from a map of optional string values,
// skipping any keys whose value pointer is nil. It is used internally by
// service methods to translate *Params structs into query strings.
type queryValues struct {
	values url.Values
}

// newQueryValues creates an empty queryValues helper.
func newQueryValues() *queryValues {
	return &queryValues{values: url.Values{}}
}

// setString adds key=v to the query if v is non-nil.
func (q *queryValues) setString(key string, v *string) {
	if v != nil {
		q.values.Set(key, *v)
	}
}

// setInt adds key=v to the query if v is non-nil.
func (q *queryValues) setInt(key string, v *int) {
	if v != nil {
		q.values.Set(key, strconv.Itoa(*v))
	}
}

// setBool adds key=v to the query if v is non-nil.
func (q *queryValues) setBool(key string, v *bool) {
	if v != nil {
		q.values.Set(key, strconv.FormatBool(*v))
	}
}

// setStringSlice adds key=a,b,c to the query if slice is non-empty.
func (q *queryValues) setStringSlice(key string, v []string) {
	if len(v) > 0 {
		q.values.Set(key, strings.Join(v, ","))
	}
}

// Values returns the underlying url.Values.
func (q *queryValues) Values() url.Values {
	return q.values
}

// doRequest performs an HTTP request against the LunarCrush API and
// decodes the JSON response body into out (if non-nil). It handles
// authentication, query string construction, JSON error mapping, and
// automatic retry with exponential backoff for HTTP 429 responses.
//
// path must be an absolute path beginning with "/" relative to the
// client's configured base URL, e.g. "/public/topic/bitcoin/v1".
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("lunarcrush: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}

	fullURL := c.baseURL + path
	if query != nil && len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	maxAttempts := c.maxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return fmt.Errorf("lunarcrush: failed to build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			lastErr = fmt.Errorf("lunarcrush: request failed: %w", err)
			return lastErr
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("lunarcrush: failed to read response body: %w", readErr)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, out); err != nil {
					return fmt.Errorf("lunarcrush: failed to decode response body: %w", err)
				}
			}
			return nil
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxAttempts {
			delay := retryDelay(resp.Header.Get("Retry-After"), c.baseDelay, attempt)
			lastErr = newAPIError(resp.StatusCode, extractMessage(respBody), respBody)

			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			continue
		}

		return newAPIError(resp.StatusCode, extractMessage(respBody), respBody)
	}

	if lastErr != nil {
		return lastErr
	}
	return newAPIError(http.StatusTooManyRequests, "rate limited", nil)
}

// extractMessage attempts to pull a human-readable message field out of a
// JSON error response body. If the body cannot be parsed, a generic
// message including a truncated body snippet is returned.
func extractMessage(body []byte) string {
	var parsed struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil {
		if parsed.Message != "" {
			return parsed.Message
		}
		if parsed.Error != "" {
			return parsed.Error
		}
	}
	if len(body) == 0 {
		return "unknown error"
	}
	snippet := string(body)
	if len(snippet) > 200 {
		snippet = snippet[:200]
	}
	return snippet
}

// retryDelay computes the delay to wait before the next retry attempt.
// It honors a Retry-After header (in seconds) if present and valid,
// otherwise it falls back to exponential backoff based on baseDelay and
// the current attempt number.
func retryDelay(retryAfterHeader string, baseDelay time.Duration, attempt int) time.Duration {
	if retryAfterHeader != "" {
		if secs, err := strconv.Atoi(strings.TrimSpace(retryAfterHeader)); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
	}
	if baseDelay <= 0 {
		baseDelay = DefaultBaseDelay
	}
	multiplier := 1 << (attempt - 1)
	return baseDelay * time.Duration(multiplier)
}
