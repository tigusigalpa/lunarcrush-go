package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestCoinsService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/coins/list/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("sort") != "galaxy_score" {
			t.Errorf("expected sort=galaxy_score, got %s", r.URL.Query().Get("sort"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"symbol":"BTC","name":"Bitcoin","galaxy_score":75.5}]}`))
	})
	defer srv.Close()

	sort := "galaxy_score"
	limit := 50
	resp, err := c.Coins.List(context.Background(), &CoinsListParams{Sort: &sort, Limit: &limit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Symbol != "BTC" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCoinsService_ListV2(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/coins/list/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"symbol":"ETH"}]}`))
	})
	defer srv.Close()

	resp, err := c.Coins.ListV2(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Symbol != "ETH" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCoinsService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/coins/bitcoin/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"symbol":"BTC","price":65000.5}}`))
	})
	defer srv.Close()

	resp, err := c.Coins.Get(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Price != 65000.5 {
		t.Errorf("expected price 65000.5, got %f", resp.Data.Price)
	}
}

func TestCoinsService_Meta(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/coins/bitcoin/meta/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"symbol":"BTC","description":"digital gold"}}`))
	})
	defer srv.Close()

	resp, err := c.Coins.Meta(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Description != "digital gold" {
		t.Errorf("unexpected description: %s", resp.Data.Description)
	}
}

func TestCoinsService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/coins/bitcoin/time-series/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("interval") != "1w" {
			t.Errorf("expected interval=1w, got %s", r.URL.Query().Get("interval"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000,"close":65000.0}]}`))
	})
	defer srv.Close()

	interval := "1w"
	resp, err := c.Coins.TimeSeries(context.Background(), "bitcoin", &TimeSeriesParams{Interval: &interval})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Close != 65000.0 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestCoinsService_Get_NotFound(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"coin not found"}`))
	})
	defer srv.Close()

	_, err := c.Coins.Get(context.Background(), "doesnotexist")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
