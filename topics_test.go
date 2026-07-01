package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestTopicsService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"topic":"bitcoin","interactions_24h":123456}}`))
	})
	defer srv.Close()

	resp, err := c.Topics.Get(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Topic != "bitcoin" || resp.Data.Interactions24h != 123456 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/time-series/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("bucket") != "hour" {
			t.Errorf("expected bucket=hour, got %s", r.URL.Query().Get("bucket"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000,"interactions":42}]}`))
	})
	defer srv.Close()

	bucket := "hour"
	resp, err := c.Topics.TimeSeries(context.Background(), "bitcoin", &TimeSeriesParams{Bucket: &bucket})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Interactions != 42 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_TimeSeriesV2(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/time-series/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000}]}`))
	})
	defer srv.Close()

	resp, err := c.Topics.TimeSeriesV2(context.Background(), "bitcoin", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_Creators(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/creators/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"creator_name":"elonmusk","creator_network":"twitter"}]}`))
	})
	defer srv.Close()

	resp, err := c.Topics.Creators(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].CreatorName != "elonmusk" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_News(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/news/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_title":"Bitcoin hits new high"}]}`))
	})
	defer srv.Close()

	resp, err := c.Topics.News(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].PostTitle != "Bitcoin hits new high" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_Posts(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/posts/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_type":"tweet"}]}`))
	})
	defer srv.Close()

	resp, err := c.Topics.Posts(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].PostType != "tweet" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_WhatsUp(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topic/bitcoin/whatsup/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"topic":"bitcoin","summary":"Bitcoin is trending"}}`))
	})
	defer srv.Close()

	resp, err := c.Topics.WhatsUp(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Summary != "Bitcoin is trending" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestTopicsService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/topics/list/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"topic":"bitcoin"},{"topic":"ethereum"}]}`))
	})
	defer srv.Close()

	resp, err := c.Topics.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
