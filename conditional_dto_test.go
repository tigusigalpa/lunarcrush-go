package lunarcrush

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func conditionalFixture(t *testing.T, name string) []byte {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", "conditional-dto", name))
	if err != nil {
		t.Fatalf("read fixture %q: %v", name, err)
	}
	return body
}

type rawResponder interface {
	RawResponse() json.RawMessage
}

func assertRawReceipt(t *testing.T, response rawResponder, fixture []byte, wantHash string) {
	t.Helper()
	receipt := response.RawResponse()
	if !bytes.Equal(receipt, fixture) {
		t.Fatal("response receipt does not match original fixture bytes")
	}
	gotHash := sha256.Sum256(receipt)
	if got := hex.EncodeToString(gotHash[:]); got != wantHash {
		t.Fatalf("response receipt hash = %s, want %s", got, wantHash)
	}
	if len(receipt) > 0 {
		receipt[0] ^= 0xff
		if bytes.Equal(receipt, response.RawResponse()) {
			t.Fatal("RawResponse returned a mutable internal receipt")
		}
	}
}

func TestConditionalDTOFixtures(t *testing.T) {
	t.Run("LC1 trending topics", func(t *testing.T) {
		fixture := conditionalFixture(t, "lc1-topics-list.json")
		c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/public/topics/list/v1"; got != want {
				t.Errorf("path = %q, want %q", got, want)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(fixture); err != nil {
				t.Errorf("write fixture: %v", err)
			}
		})
		defer srv.Close()

		response, err := c.Topics.List(context.Background())
		if err != nil {
			t.Fatalf("Topics.List: %v", err)
		}
		if len(response.Data) != 1 {
			t.Fatalf("topic count = %d, want 1", len(response.Data))
		}
		item := response.Data[0]
		if item.TopicRank1hPrevious == nil || *item.TopicRank1hPrevious != 2 {
			t.Fatalf("topic_rank_1h_previous = %v, want 2", item.TopicRank1hPrevious)
		}
		if item.TopicRank24hPrevious == nil || *item.TopicRank24hPrevious != 3 {
			t.Fatalf("topic_rank_24h_previous = %v, want 3", item.TopicRank24hPrevious)
		}
		if item.NumPosts == nil || *item.NumPosts != 0 {
			t.Fatalf("num_posts = %v, want present zero", item.NumPosts)
		}
		assertRawReceipt(t, response, fixture, "5df702dc1415e51610931efdc4caae8e8fd5adcf1610ab213f7f34b2b695d39e")
	})

	t.Run("LC2 topic summary", func(t *testing.T) {
		fixture := conditionalFixture(t, "lc2-topic-summary.json")
		c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(fixture); err != nil {
				t.Errorf("write fixture: %v", err)
			}
		})
		defer srv.Close()

		response, err := c.Topics.Get(context.Background(), "bitcoin")
		if err != nil {
			t.Fatalf("Topics.Get: %v", err)
		}
		if got, want := response.Data.RelatedTopics, []string{"crypto"}; !equalStrings(got, want) {
			t.Fatalf("related_topics = %#v, want %#v", got, want)
		}
		if got := response.Data.TypesSentimentDetail["tweet"]["positive"]; got != 25741 {
			t.Fatalf("tweet positive sentiment detail = %v, want 25741", got)
		}
		if !bytes.Contains(response.RawResponse(), []byte(`"unknown_network_key"`)) {
			t.Fatal("raw receipt did not preserve an unknown provider field")
		}
		assertRawReceipt(t, response, fixture, "efc6777fd53ac63d138f72db0ed8a1127011ca9959bffcd769812af8404e8b53")
	})

	t.Run("LC3 coin snapshot", func(t *testing.T) {
		fixture := conditionalFixture(t, "lc3-coin-snapshot.json")
		c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(fixture); err != nil {
				t.Errorf("write fixture: %v", err)
			}
		})
		defer srv.Close()

		response, err := c.Coins.Get(context.Background(), "ETH")
		if err != nil {
			t.Fatalf("Coins.Get: %v", err)
		}
		if response.Data.MaxSupply != nil {
			t.Fatalf("max_supply = %v, want nil", *response.Data.MaxSupply)
		}
		if response.Data.Price != 0 {
			t.Fatalf("price = %v, want zero", response.Data.Price)
		}
		assertRawReceipt(t, response, fixture, "3523ae157fa9c643f52a74f1dea01ab05959ad96a6052c0acada0911d015508e")
	})

	t.Run("LC4 topic time series", func(t *testing.T) {
		fixture := conditionalFixture(t, "lc4-topic-time-series.json")
		c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if got, want := r.URL.Path, "/public/topic/bitcoin/time-series/v2"; got != want {
				t.Errorf("path = %q, want %q", got, want)
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(fixture); err != nil {
				t.Errorf("write fixture: %v", err)
			}
		})
		defer srv.Close()

		response, err := c.Topics.TimeSeriesV2(context.Background(), "bitcoin", nil)
		if err != nil {
			t.Fatalf("Topics.TimeSeriesV2: %v", err)
		}
		if len(response.Data) != 1 {
			t.Fatalf("time-series point count = %d, want 1", len(response.Data))
		}
		point := response.Data[0]
		if point.Time != 1783209600 {
			t.Fatalf("row time = %d, want Unix seconds 1783209600", point.Time)
		}
		if point.PostsActive == nil || *point.PostsActive != 0 {
			t.Fatalf("posts_active = %v, want present zero", point.PostsActive)
		}
		if point.PostsCreated == nil || *point.PostsCreated != 1070 {
			t.Fatalf("posts_created = %v, want 1070", point.PostsCreated)
		}
		assertRawReceipt(t, response, fixture, "6c7363e5cee50f2f5df5c12894752cbff0e95a1a9d46476832e58f94ffc595a7")
	})
}

func TestNullableFieldsPreserveNullAndZero(t *testing.T) {
	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":[{"post_title":null,"post_sentiment":null,"post_created":null},{},{"post_title":"","post_sentiment":0,"post_created":0}]}`)
	})
	defer srv.Close()

	response, err := c.Topics.Posts(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("Topics.Posts: %v", err)
	}
	if len(response.Data) != 3 {
		t.Fatalf("post count = %d, want 3", len(response.Data))
	}
	if response.Data[0].PostTitle != nil || response.Data[0].PostSentiment != nil || response.Data[0].CreatedTime != nil {
		t.Fatal("null post values must remain nil")
	}
	if response.Data[1].PostTitle != nil || response.Data[1].PostSentiment != nil || response.Data[1].CreatedTime != nil {
		t.Fatal("absent post values must remain nil in the convenience DTO")
	}
	zero := response.Data[2]
	if zero.PostTitle == nil || *zero.PostTitle != "" || zero.PostSentiment == nil || *zero.PostSentiment != 0 || zero.CreatedTime == nil || *zero.CreatedTime != 0 {
		t.Fatal("present empty and zero post values must not be coerced to nil")
	}

	var raw struct {
		Data []map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(response.RawResponse(), &raw); err != nil {
		t.Fatalf("decode raw receipt: %v", err)
	}
	if got, ok := raw.Data[0]["post_title"]; !ok || string(got) != "null" {
		t.Fatalf("raw receipt lost explicit null post_title: %q", got)
	}
	if _, ok := raw.Data[1]["post_title"]; ok {
		t.Fatal("raw receipt did not preserve an absent post_title")
	}
}

func TestCoinMaxSupplyPresentZeroAndNativeTimeSeriesPath(t *testing.T) {
	var coin Coin
	if err := json.Unmarshal([]byte(`{"max_supply":0}`), &coin); err != nil {
		t.Fatalf("decode coin: %v", err)
	}
	if coin.MaxSupply == nil || *coin.MaxSupply != 0 {
		t.Fatalf("max_supply = %v, want present zero", coin.MaxSupply)
	}

	c, srv := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/public/coins/ETH/time-series/v2"; got != want {
			t.Errorf("path = %q, want %q", got, want)
		}
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, `{"data":[]}`)
	})
	defer srv.Close()
	if _, err := c.Coins.TimeSeries(context.Background(), "ETH", nil); err != nil {
		t.Fatalf("Coins.TimeSeries: %v", err)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
