package lunarcrush

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestServicesReturnAPIError(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		writeResponse(t, w, `{"error":"service unavailable"}`)
	})
	defer srv.Close()

	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
	}{
		{"AI.Topic", func() error { _, err := c.AI.Topic(ctx, "bitcoin"); return err }},
		{"AI.Creator", func() error { _, err := c.AI.Creator(ctx, "twitter", "creator"); return err }},
		{"Categories.List", func() error { _, err := c.Categories.List(ctx); return err }},
		{"Categories.Get", func() error { _, err := c.Categories.Get(ctx, "defi"); return err }},
		{"Categories.Creators", func() error { _, err := c.Categories.Creators(ctx, "defi"); return err }},
		{"Categories.News", func() error { _, err := c.Categories.News(ctx, "defi"); return err }},
		{"Categories.Posts", func() error { _, err := c.Categories.Posts(ctx, "defi"); return err }},
		{"Categories.TimeSeries", func() error { _, err := c.Categories.TimeSeries(ctx, "defi", nil); return err }},
		{"Categories.Topics", func() error { _, err := c.Categories.Topics(ctx, "defi"); return err }},
		{"Coins.List", func() error { _, err := c.Coins.List(ctx, nil); return err }},
		{"Coins.ListV2", func() error { _, err := c.Coins.ListV2(ctx, nil); return err }},
		{"Coins.Get", func() error { _, err := c.Coins.Get(ctx, "BTC"); return err }},
		{"Coins.Meta", func() error { _, err := c.Coins.Meta(ctx, "BTC"); return err }},
		{"Coins.TimeSeries", func() error { _, err := c.Coins.TimeSeries(ctx, "BTC", nil); return err }},
		{"Creators.Get", func() error { _, err := c.Creators.Get(ctx, "twitter", "creator"); return err }},
		{"Creators.Posts", func() error { _, err := c.Creators.Posts(ctx, "twitter", "creator"); return err }},
		{"Creators.TimeSeries", func() error { _, err := c.Creators.TimeSeries(ctx, "twitter", "creator", nil); return err }},
		{"Creators.List", func() error { _, err := c.Creators.List(ctx); return err }},
		{"Posts.List", func() error { _, err := c.Posts.List(ctx, nil); return err }},
		{"Posts.TimeSeries", func() error { _, err := c.Posts.TimeSeries(ctx, nil); return err }},
		{"Searches.Create", func() error {
			_, err := c.Searches.Create(ctx, &SearchCreateParams{Name: "search", SearchJSON: `{}`})
			return err
		}},
		{"Searches.List", func() error { _, err := c.Searches.List(ctx); return err }},
		{"Searches.Search", func() error { _, err := c.Searches.Search(ctx, `{}`); return err }},
		{"Searches.Get", func() error { _, err := c.Searches.Get(ctx, "search"); return err }},
		{"Searches.Update", func() error { _, err := c.Searches.Update(ctx, "search", nil); return err }},
		{"Searches.Delete", func() error { _, err := c.Searches.Delete(ctx, "search"); return err }},
		{"Stocks.List", func() error { _, err := c.Stocks.List(ctx, nil); return err }},
		{"Stocks.ListV2", func() error { _, err := c.Stocks.ListV2(ctx, nil); return err }},
		{"Stocks.Get", func() error { _, err := c.Stocks.Get(ctx, "NVDA"); return err }},
		{"Stocks.TimeSeries", func() error { _, err := c.Stocks.TimeSeries(ctx, "NVDA", nil); return err }},
		{"System.Changes", func() error { _, err := c.System.Changes(ctx); return err }},
		{"Topics.Get", func() error { _, err := c.Topics.Get(ctx, "bitcoin"); return err }},
		{"Topics.TimeSeries", func() error { _, err := c.Topics.TimeSeries(ctx, "bitcoin", nil); return err }},
		{"Topics.TimeSeriesV2", func() error { _, err := c.Topics.TimeSeriesV2(ctx, "bitcoin", nil); return err }},
		{"Topics.Creators", func() error { _, err := c.Topics.Creators(ctx, "bitcoin"); return err }},
		{"Topics.News", func() error { _, err := c.Topics.News(ctx, "bitcoin"); return err }},
		{"Topics.Posts", func() error { _, err := c.Topics.Posts(ctx, "bitcoin"); return err }},
		{"Topics.WhatsUp", func() error { _, err := c.Topics.WhatsUp(ctx, "bitcoin"); return err }},
		{"Topics.List", func() error { _, err := c.Topics.List(ctx); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiErr *APIError
			if err := tt.call(); !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusInternalServerError {
				t.Fatalf("expected API error with status 500, got %v", err)
			}
		})
	}
}

func TestRequestHelpers(t *testing.T) {
	t.Run("extract message", func(t *testing.T) {
		cases := []struct {
			body string
			want string
		}{
			{`{"message":"message"}`, "message"},
			{`{"error":"error"}`, "error"},
			{"not json", "not json"},
			{"", "unknown error"},
		}
		for _, tt := range cases {
			if got := extractMessage([]byte(tt.body)); got != tt.want {
				t.Errorf("extractMessage(%q) = %q, want %q", tt.body, got, tt.want)
			}
		}
		if got := len(extractMessage(make([]byte, 201))); got != 200 {
			t.Errorf("expected a 200-byte truncated message, got %d", got)
		}
	})

	t.Run("retry delay", func(t *testing.T) {
		if got := retryDelay("2", time.Millisecond, 1); got != 2*time.Second {
			t.Errorf("expected Retry-After delay 2s, got %s", got)
		}
		if got := retryDelay("invalid", time.Millisecond, 3); got != 4*time.Millisecond {
			t.Errorf("expected exponential delay 4ms, got %s", got)
		}
		if got := retryDelay("", 0, 1); got != DefaultBaseDelay {
			t.Errorf("expected default delay %s, got %s", DefaultBaseDelay, got)
		}
	})

	t.Run("wrapped API error", func(t *testing.T) {
		wrapped := newAPIError(http.StatusUnauthorized, "unauthorized", nil).(*wrappedAPIError)
		if got := wrapped.Unwrap(); !errors.Is(got, ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", got)
		}
		var apiErr *APIError
		if !wrapped.As(&apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected APIError from As, got %v", apiErr)
		}
		var pathErr *os.PathError
		if wrapped.As(&pathErr) {
			t.Error("did not expect an unrelated type to match As")
		}
	})
}

func TestRequestFailurePaths(t *testing.T) {
	c := NewClient("key")
	if err := c.doRequest(context.Background(), http.MethodPost, "/test", nil, make(chan int), nil); err == nil {
		t.Fatal("expected marshal error")
	}

	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, "not json")
	})
	defer srv.Close()
	var out map[string]interface{}
	if err := c.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, &out); err == nil {
		t.Fatal("expected decode error")
	}

	transportErr := errors.New("transport failure")
	c = NewClient("key", WithHTTPClient(&http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, transportErr
	})}))
	if err := c.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, nil); !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestPostsTimeSeriesParamsToQuery(t *testing.T) {
	topic, category, bucket, interval := "bitcoin", "defi", "hour", "1w"
	start, end := 10, 20
	params := &PostsTimeSeriesParams{
		Topic:    &topic,
		Category: &category,
		Bucket:   &bucket,
		Interval: &interval,
		Start:    &start,
		End:      &end,
	}
	values := params.toQuery().Values()
	if got, want := values.Encode(), "bucket=hour&category=defi&end=20&interval=1w&start=10&topic=bitcoin"; got != want {
		t.Errorf("unexpected query values: got %q, want %q", got, want)
	}
}
