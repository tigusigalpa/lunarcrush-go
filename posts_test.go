package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestPostsService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/posts/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("topic") != "bitcoin" {
			t.Errorf("expected topic=bitcoin, got %s", r.URL.Query().Get("topic"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_title":"BTC to the moon"}]}`))
	})
	defer srv.Close()

	topic := "bitcoin"
	resp, err := c.Posts.List(context.Background(), &PostsListParams{Topic: &topic})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].PostTitle != "BTC to the moon" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestPostsService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/posts/time-series/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000,"posts_created":50}]}`))
	})
	defer srv.Close()

	resp, err := c.Posts.TimeSeries(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Posts != 50 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
