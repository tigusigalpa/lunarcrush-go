package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestStocksService_List(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/stocks/list/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"symbol":"NVDA","price":120.5}]}`))
	})
	defer srv.Close()

	resp, err := c.Stocks.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Symbol != "NVDA" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestStocksService_ListV2(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/stocks/list/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"symbol":"AAPL"}]}`))
	})
	defer srv.Close()

	sort := "market_cap"
	resp, err := c.Stocks.ListV2(context.Background(), &StocksListParams{Sort: &sort})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Symbol != "AAPL" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestStocksService_Get(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/stocks/NVDA/v1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"symbol":"NVDA","galaxy_score":80}}`))
	})
	defer srv.Close()

	resp, err := c.Stocks.Get(context.Background(), "NVDA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.GalaxyScore != 80 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestStocksService_TimeSeries(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/stocks/NVDA/time-series/v2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"time":1700000000,"close":120.5}]}`))
	})
	defer srv.Close()

	resp, err := c.Stocks.TimeSeries(context.Background(), "NVDA", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Close != 120.5 {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestStocksService_Get_ServerError(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal error"}`))
	})
	defer srv.Close()

	_, err := c.Stocks.Get(context.Background(), "NVDA")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if ae, ok := err.(*APIError); ok {
		apiErr = ae
	}
	if apiErr == nil || apiErr.StatusCode != 500 {
		t.Errorf("expected APIError with status 500, got %v", err)
	}
}
