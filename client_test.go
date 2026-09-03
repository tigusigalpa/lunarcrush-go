package lunarcrush

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient spins up an httptest.Server using handler and returns a
// *Client configured to talk to it, along with the server for cleanup.
func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	allOpts := append([]Option{WithBaseURL(srv.URL)}, opts...)
	c := NewClient("test-api-key", allOpts...)
	return c, srv
}

// writeResponse writes an HTTP response body and reports an unexpected write
// failure to the active test.
func writeResponse(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()
	if _, err := w.Write([]byte(body)); err != nil {
		t.Errorf("failed to write test response: %v", err)
	}
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("key")
	if c.baseURL != DefaultBaseURL {
		t.Errorf("expected baseURL %s, got %s", DefaultBaseURL, c.baseURL)
	}
	if c.httpClient.Timeout != DefaultTimeout {
		t.Errorf("expected timeout %s, got %s", DefaultTimeout, c.httpClient.Timeout)
	}
	if c.maxAttempts != DefaultMaxAttempts {
		t.Errorf("expected maxAttempts %d, got %d", DefaultMaxAttempts, c.maxAttempts)
	}
	if c.Coins == nil || c.Topics == nil || c.Stocks == nil || c.Creators == nil ||
		c.Categories == nil || c.Posts == nil || c.Searches == nil || c.System == nil || c.AI == nil {
		t.Error("expected all services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customClient := &http.Client{}
	c := NewClient("key",
		WithBaseURL("https://example.com"),
		WithHTTPClient(customClient),
		WithTimeout(5*time.Second),
		WithRetry(3, 100*time.Millisecond),
	)
	if c.baseURL != "https://example.com" {
		t.Errorf("expected custom baseURL, got %s", c.baseURL)
	}
	if c.httpClient != customClient {
		t.Error("expected custom http client to be used")
	}
	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %s", c.httpClient.Timeout)
	}
	if c.maxAttempts != 3 {
		t.Errorf("expected maxAttempts 3, got %d", c.maxAttempts)
	}
	if c.baseDelay != 100*time.Millisecond {
		t.Errorf("expected baseDelay 100ms, got %s", c.baseDelay)
	}
}

func TestWithBaseURL_TrailingSlash(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/public/topic/bitcoin/v1"; got != want {
			t.Errorf("expected path %q, got %q", want, got)
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{"topic":"bitcoin"}}`)
	}))
	defer srv.Close()

	c := NewClient("test-api-key", WithBaseURL(srv.URL+"/"))
	if _, err := c.Topics.Get(context.Background(), "bitcoin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPathParametersAreEscaped(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.EscapedPath(), "/public/topic/a%2Fb%3Fc%23d%20e/v1"; got != want {
			t.Errorf("expected escaped path %q, got %q", want, got)
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{}}`)
	})
	defer srv.Close()

	if _, err := c.Topics.Get(context.Background(), "a/b?c#d e"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_RetryPreservesRequestBody(t *testing.T) {
	attempts := 0
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		defer r.Body.Close()
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body on attempt %d: %v", attempts, err)
		}
		if got, want := body["name"], "saved search"; got != want {
			t.Errorf("expected body name %q, got %q", want, got)
		}
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{}}`)
	}, WithRetry(2, time.Millisecond))
	defer srv.Close()

	var out map[string]interface{}
	if err := c.doRequest(context.Background(), http.MethodPost, "/test", nil, map[string]string{"name": "saved search"}, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestDoRequest_Retry429(t *testing.T) {
	var attempts int
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			writeResponse(t, w, `{"message":"rate limited"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{"topic":"bitcoin"}}`)
	}, WithRetry(5, 10*time.Millisecond))
	defer srv.Close()

	resp, err := c.Topics.Get(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("expected no error after retries, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if resp.Data.Topic != "bitcoin" {
		t.Errorf("expected topic bitcoin, got %s", resp.Data.Topic)
	}
}

func TestDoRequest_RateLimitedExhausted(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeResponse(t, w, `{"message":"rate limited"}`)
	}, WithRetry(2, 5*time.Millisecond))
	defer srv.Close()

	_, err := c.Topics.Get(context.Background(), "bitcoin")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrRateLimited) {
		t.Errorf("expected ErrRateLimited, got %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Errorf("expected *APIError, got %T", err)
	} else if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", apiErr.StatusCode)
	}
}

func TestDoRequest_Unauthorized(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		writeResponse(t, w, `{"message":"invalid api key"}`)
	})
	defer srv.Close()

	_, err := c.Topics.Get(context.Background(), "bitcoin")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeResponse(t, w, `{"message":"not found"}`)
	})
	defer srv.Close()

	_, err := c.Topics.Get(context.Background(), "unknown")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDoRequest_ContextCancellation(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{}}`)
	})
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := c.Topics.Get(ctx, "bitcoin")
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestDoRequest_AuthorizationHeader(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-api-key" {
			t.Errorf("expected Authorization header 'Bearer test-api-key', got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":{}}`)
	})
	defer srv.Close()

	if _, err := c.Topics.Get(context.Background(), "bitcoin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 500, Message: "server error"}
	want := "lunarcrush: api error: status=500 message=server error"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}
