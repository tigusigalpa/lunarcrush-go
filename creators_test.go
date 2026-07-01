package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestCreatorsService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/creator/twitter/elonmusk/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"creator_name":"elonmusk","creator_network":"twitter","creator_followers":150000000}}`))
	})
	defer srv.Close()

	resp, err := c.Creators.Get(context.Background(), "twitter", "elonmusk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.CreatorName != "elonmusk" || resp.Data.FollowerCount != 150000000 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCreatorsService_Posts(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/creator/twitter/elonmusk/posts/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_title":"hello"}]}`))
	})
	defer srv.Close()

	resp, err := c.Creators.Posts(context.Background(), "twitter", "elonmusk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].PostTitle != "hello" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCreatorsService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/creator/twitter/elonmusk/time-series/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000,"interactions":99}]}`))
	})
	defer srv.Close()

	resp, err := c.Creators.TimeSeries(context.Background(), "twitter", "elonmusk", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Interactions != 99 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCreatorsService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/creators/list/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"creator_name":"a"},{"creator_name":"b"}]}`))
	})
	defer srv.Close()

	resp, err := c.Creators.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
