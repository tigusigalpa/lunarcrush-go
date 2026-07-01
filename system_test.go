package lunarcrush

import (
	"context"
	"net/http"
	"testing"
)

func TestSystemService_Changes(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/public/system/changes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"category":"api","change":"added endpoint","timestamp":1700000000}]}`))
	})
	defer srv.Close()

	resp, err := c.System.Changes(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 || resp.Data[0].Change != "added endpoint" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
}

func TestSystemService_Changes_Error(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"boom"}`))
	})
	defer srv.Close()

	_, err := c.System.Changes(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
