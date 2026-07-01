package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestCategoriesService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/categories/list/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"category":"defi"}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Category != "defi" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"category":"defi","sentiment":0.8}}`))
	})
	defer srv.Close()

	resp, err := c.Categories.Get(context.Background(), "defi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Sentiment != 0.8 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_Creators(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/creators/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"creator_name":"c1"}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.Creators(context.Background(), "defi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_News(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/news/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_title":"DeFi news"}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.News(context.Background(), "defi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].PostTitle != "DeFi news" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_Posts(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/posts/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"post_type":"tweet"}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.Posts(context.Background(), "defi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/time-series/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.TimeSeries(context.Background(), "defi", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCategoriesService_Topics(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/category/defi/topics/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"topic":"uniswap"}]}`))
	})
	defer srv.Close()

	resp, err := c.Categories.Topics(context.Background(), "defi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Topic != "uniswap" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
