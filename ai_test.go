package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestAIService_Topic(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/ai/topic/bitcoin" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"subject":"bitcoin","summary":"Bullish sentiment overall."}}`))
	})
	defer srv.Close()

	resp, err := c.AI.Topic(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Summary != "Bullish sentiment overall." {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestAIService_Creator(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/ai/creator/twitter/elonmusk" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"subject":"elonmusk","summary":"High engagement creator."}}`))
	})
	defer srv.Close()

	resp, err := c.AI.Creator(context.Background(), "twitter", "elonmusk")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Summary != "High engagement creator." {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}
