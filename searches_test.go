package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestSearchesService_Create(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/create" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("name") != "my-search" {
			t.Errorf("expected name=my-search, got %s", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("search_json") != `{"q":"bitcoin"}` {
			t.Errorf("unexpected search_json: %s", r.URL.Query().Get("search_json"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"slug":"my-search-abc","name":"my-search"}}`))
	})
	defer srv.Close()

	resp, err := c.Searches.Create(context.Background(), &SearchCreateParams{
		Name:       "my-search",
		SearchJSON: `{"q":"bitcoin"}`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Slug != "my-search-abc" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSearchesService_Create_MissingRequiredFields(t *testing.T) {
	c := NewClient("key")
	_, err := c.Searches.Create(context.Background(), &SearchCreateParams{Name: ""})
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestSearchesService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/list" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"slug":"s1"},{"slug":"s2"}]}`))
	})
	defer srv.Close()

	resp, err := c.Searches.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSearchesService_Search(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"slug":"ad-hoc"}}`))
	})
	defer srv.Close()

	resp, err := c.Searches.Search(context.Background(), `{"q":"eth"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Slug != "ad-hoc" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSearchesService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/my-search-abc" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"slug":"my-search-abc"}}`))
	})
	defer srv.Close()

	resp, err := c.Searches.Get(context.Background(), "my-search-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Slug != "my-search-abc" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSearchesService_Update(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/my-search-abc/update" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("priority") != "5" {
			t.Errorf("expected priority=5, got %s", r.URL.Query().Get("priority"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"slug":"my-search-abc","priority":5}}`))
	})
	defer srv.Close()

	priority := 5
	resp, err := c.Searches.Update(context.Background(), "my-search-abc", &SearchUpdateParams{Priority: &priority})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Priority != 5 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSearchesService_Delete(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/searches/my-search-abc/delete" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"deleted":true}}`))
	})
	defer srv.Close()

	resp, err := c.Searches.Delete(context.Background(), "my-search-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted, ok := resp.Data["deleted"].(bool); !ok || !deleted {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
