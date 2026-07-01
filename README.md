# lunarcrush-go

**A fast, dependency-free Go SDK for the LunarCrush API v4.** Pull real-time and historical social & market data for crypto, stocks, topics, categories, and creators — with idiomatic Go types, functional-options configuration, context-aware requests, and automatic retry on rate limits.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/lunarcrush-go.svg)](https://pkg.go.dev/github.com/tigusigalpa/lunarcrush-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](.)

## Installation

```bash
go get github.com/tigusigalpa/lunarcrush-go
```

Requires Go 1.21+. No third-party dependencies — built entirely on the standard library (`net/http`, `encoding/json`, `context`, `time`).

## Client Initialization

```go
package main

import (
    "time"

    lunarcrush "github.com/tigusigalpa/lunarcrush-go"
)

func main() {
    client := lunarcrush.NewClient("YOUR_API_KEY",
        lunarcrush.WithBaseURL("https://lunarcrush.com/api4"), // optional, this is the default
        lunarcrush.WithTimeout(15*time.Second),
        lunarcrush.WithRetry(3, time.Second), // retry up to 3 times on HTTP 429, exponential backoff starting at 1s
        lunarcrush.WithHTTPClient(&http.Client{}), // optional custom *http.Client
    )
    _ = client
}
```

| Option | Description |
|---|---|
| `WithBaseURL(url string)` | Overrides the API base URL (default `https://lunarcrush.com/api4`). Useful for testing. |
| `WithHTTPClient(client *http.Client)` | Supplies a custom HTTP client (proxies, transports, etc). |
| `WithTimeout(d time.Duration)` | Sets the HTTP client timeout (default 30s). |
| `WithRetry(maxAttempts int, baseDelay time.Duration)` | Enables automatic retry with exponential backoff on HTTP 429, honoring `Retry-After`. |

## Quick Start

```go
ctx := context.Background()

// 1. Get a 24h topic summary
topic, err := client.Topics.Get(ctx, "bitcoin")

// 2. Get topic time-series data
bucket, interval := "hour", "1w"
series, err := client.Topics.TimeSeries(ctx, "bitcoin", &lunarcrush.TimeSeriesParams{
    Bucket:   &bucket,
    Interval: &interval,
})

// 3. List top coins by galaxy score
sort := "galaxy_score"
limit := 50
coins, err := client.Coins.List(ctx, &lunarcrush.CoinsListParams{
    Sort:  &sort,
    Limit: &limit,
    Desc:  boolPtr(true),
})

// 4. Get a single coin
coin, err := client.Coins.Get(ctx, "bitcoin")

// 5. Get a creator profile
creator, err := client.Creators.Get(ctx, "twitter", "elonmusk")

// 6. Get stock data
stock, err := client.Stocks.Get(ctx, "NVDA")

// 7. Create a saved search
search, err := client.Searches.Create(ctx, &lunarcrush.SearchCreateParams{
    Name:       "my-search",
    SearchJSON: `{"q":"bitcoin"}`,
})

func boolPtr(b bool) *bool { return &b }
```

## API Reference

All methods accept `context.Context` as their first argument and are safe for concurrent use.

### Topics — `client.Topics`

| Method | Endpoint | Description |
|---|---|---|
| `Get(ctx, topic)` | `GET /public/topic/:topic/v1` | 24h social summary |
| `TimeSeries(ctx, topic, params)` | `GET /public/topic/:topic/time-series/v1` | Historical time-series |
| `TimeSeriesV2(ctx, topic, params)` | `GET /public/topic/:topic/time-series/v2` | Historical time-series (v2) |
| `Creators(ctx, topic)` | `GET /public/topic/:topic/creators/v1` | Top creators for topic |
| `News(ctx, topic)` | `GET /public/topic/:topic/news/v1` | News articles |
| `Posts(ctx, topic)` | `GET /public/topic/:topic/posts/v1` | Social posts |
| `WhatsUp(ctx, topic)` | `GET /public/topic/:topic/whatsup/v1` | AI-generated summary |
| `List(ctx)` | `GET /public/topics/list/v1` | List all topics |

### Categories — `client.Categories`

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx)` | `GET /public/categories/list/v1` | List all categories |
| `Get(ctx, category)` | `GET /public/category/:category/v1` | Category summary |
| `Creators(ctx, category)` | `GET /public/category/:category/creators/v1` | Top creators |
| `News(ctx, category)` | `GET /public/category/:category/news/v1` | News articles |
| `Posts(ctx, category)` | `GET /public/category/:category/posts/v1` | Social posts |
| `TimeSeries(ctx, category, params)` | `GET /public/category/:category/time-series/v1` | Historical time-series |
| `Topics(ctx, category)` | `GET /public/category/:category/topics/v1` | Topics in category |

### Creators — `client.Creators`

| Method | Endpoint | Description |
|---|---|---|
| `Get(ctx, network, id)` | `GET /public/creator/:network/:id/v1` | Creator profile |
| `Posts(ctx, network, id)` | `GET /public/creator/:network/:id/posts/v1` | Creator's posts |
| `TimeSeries(ctx, network, id, params)` | `GET /public/creator/:network/:id/time-series/v1` | Historical time-series |
| `List(ctx)` | `GET /public/creators/list/v1` | List top creators |

Supported `network` values: `twitter`, `youtube`, `instagram`, `reddit`, `tiktok`.

### Posts — `client.Posts`

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx, params)` | `GET /public/posts/v1` | List posts |
| `TimeSeries(ctx, params)` | `GET /public/posts/time-series/v1` | Post activity time-series |

### Coins — `client.Coins`

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx, params)` | `GET /public/coins/list/v1` | List coins (v1) |
| `ListV2(ctx, params)` | `GET /public/coins/list/v2` | List coins (v2) |
| `Get(ctx, coin)` | `GET /public/coins/:coin/v1` | Coin details |
| `Meta(ctx, coin)` | `GET /public/coins/:coin/meta/v1` | Coin metadata |
| `TimeSeries(ctx, coin, params)` | `GET /public/coins/:coin/time-series/v2` | Historical time-series |

### Stocks — `client.Stocks`

| Method | Endpoint | Description |
|---|---|---|
| `List(ctx, params)` | `GET /public/stocks/list/v1` | List stocks (v1) |
| `ListV2(ctx, params)` | `GET /public/stocks/list/v2` | List stocks (v2) |
| `Get(ctx, stock)` | `GET /public/stocks/:stock/v1` | Stock details |
| `TimeSeries(ctx, stock, params)` | `GET /public/stocks/:stock/time-series/v2` | Historical time-series |

### Searches — `client.Searches`

| Method | Endpoint | Description |
|---|---|---|
| `Create(ctx, params)` | `GET /public/searches/create` | Create a saved search |
| `List(ctx)` | `GET /public/searches/list` | List saved searches |
| `Search(ctx, searchJSON)` | `GET /public/searches/search` | Ad-hoc search |
| `Get(ctx, slug)` | `GET /public/searches/:slug` | Get a saved search |
| `Update(ctx, slug, params)` | `GET /public/searches/:slug/update` | Update a saved search |
| `Delete(ctx, slug)` | `GET /public/searches/:slug/delete` | Delete a saved search |

### AI — `client.AI`

| Method | Endpoint | Description |
|---|---|---|
| `Topic(ctx, topic)` | `GET /public/ai/topic/:topic` | AI summary for topic |
| `Creator(ctx, network, id)` | `GET /public/ai/creator/:network/:id` | AI summary for creator |

### System — `client.System`

| Method | Endpoint | Description |
|---|---|---|
| `Changes(ctx)` | `GET /public/system/changes` | Recent API/platform changes |

## Rate Limits

| Plan | Requests/min | Requests/day |
|---|---|---|
| Hobby | 4 | 100 |
| Individual | 10 | 2,000 |
| Builder | 100 | 20,000 |
| Scale | 500 | 100,000 |

Use `WithRetry` to automatically retry HTTP 429 responses with exponential backoff, honoring the `Retry-After` header when present.

## Error Handling

Every non-2xx response is returned as an `*lunarcrush.APIError`:

```go
type APIError struct {
    StatusCode int
    Message    string
    RawBody    []byte
}
```

Sentinel errors are provided for common cases and work with `errors.Is`:

```go
topic, err := client.Topics.Get(ctx, "unknown-topic")
if err != nil {
    switch {
    case errors.Is(err, lunarcrush.ErrNotFound):
        // handle 404
    case errors.Is(err, lunarcrush.ErrUnauthorized):
        // handle 401 — check your API key
    case errors.Is(err, lunarcrush.ErrRateLimited):
        // handle 429 — retries exhausted
    default:
        var apiErr *lunarcrush.APIError
        if errors.As(err, &apiErr) {
            fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

## Context & Concurrency

Every method accepts `context.Context` as its first argument and respects cancellation/deadlines. The `Client` and all services are safe for concurrent use across goroutines:

```go
var wg sync.WaitGroup
for _, sym := range []string{"bitcoin", "ethereum", "solana"} {
    wg.Add(1)
    go func(symbol string) {
        defer wg.Done()
        coin, err := client.Coins.Get(ctx, symbol)
        _ = coin
        _ = err
    }(sym)
}
wg.Wait()
```

## Testing

The SDK ships with a full test suite using `net/http/httptest`, covering successful deserialization, HTTP 429 retry behavior, context cancellation, and error mapping:

```bash
go test ./...
```

## Contributing

Issues and pull requests are welcome at [github.com/tigusigalpa/lunarcrush-go](https://github.com/tigusigalpa/lunarcrush-go).

1. Fork the repository.
2. Create a feature branch.
3. Add tests for any new behavior.
4. Run `go vet ./...` and `go test ./...` before submitting a PR.

## License

MIT © [Igor Sazonov](https://github.com/tigusigalpa)
