# lunarcrush-go

![LunarCrush Golang SDK](https://i.postimg.cc/0QdLHxk1/lunarcrush-go-banner.jpg)

> **Social + market intelligence, idiomatically written in Go.**

Welcome to `lunarcrush-go`, a friendly and fast Go SDK for the [LunarCrush API v4](https://lunarcrush.com/api4). Whether
you are building a trading dashboard, tracking crypto sentiment, or analyzing how a stock is trending on social media,
this library gives you clean, typed, and production-ready access to every LunarCrush endpoint — without pulling in a
single third-party dependency.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/lunarcrush-go.svg)](https://pkg.go.dev/github.com/tigusigalpa/lunarcrush-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](.)

## What you get

- **Complete API coverage** — Coins, Stocks, Topics, Categories, Creators, Posts, Searches, AI summaries, and System
  changes.
- **Truly zero dependencies** — only the Go standard library (`net/http`, `encoding/json`, `context`, `time`).
- **Functional options** — configure the client the idiomatic Go way.
- **Context-aware & concurrent** — every method accepts `context.Context` and the client is safe for use across
  goroutines.
- **Built-in resilience** — automatic retry with exponential backoff on HTTP 429, respecting `Retry-After` when
  LunarCrush tells you to wait.
- **Friendly errors** — sentinel errors for `401`, `404`, and `429`, plus detailed `APIError` values for everything
  else.
- **Fully tested** — every service has its own `httptest`-based test suite.

## Installation

Drop it into your project with one command:

```bash
go get github.com/tigusigalpa/lunarcrush-go
```

Requirements: Go 1.21 or newer. No `go.sum` bloat, no dependency tree drama — just the standard library doing what it
does best.

## Getting started in 60 seconds

Here is a tiny, complete program you can run right away. Just replace `YOUR_API_KEY` with your real LunarCrush key:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    lunarcrush "github.com/tigusigalpa/lunarcrush-go"
)

func main() {
    ctx := context.Background()

    client := lunarcrush.NewClient("YOUR_API_KEY",
        lunarcrush.WithTimeout(15*time.Second),
        lunarcrush.WithRetry(3, time.Second), // 3 attempts, 1s initial backoff
    )

    // 24-hour social summary for Bitcoin
    topic, err := client.Topics.Get(ctx, "bitcoin")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bitcoin interactions (24h): %.0f\n", topic.Data.Interactions24h)

    // Top 10 coins by galaxy score
    sort := "galaxy_score"
    limit := 10
    coins, err := client.Coins.List(ctx, &lunarcrush.CoinsListParams{
        Sort:  &sort,
        Limit: &limit,
        Desc:  ptr(true),
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, coin := range coins.Data {
        fmt.Printf("%s — galaxy score %.1f\n", coin.Symbol, coin.GalaxyScore)
    }
}

func ptr[T any](v T) *T { return &v }
```

Compile and run it:

```bash
go run main.go
```

That is it. You are now talking to LunarCrush from Go.

## Client configuration

`NewClient` uses the functional options pattern. Mix and match only what you need:

```go
client := lunarcrush.NewClient("YOUR_API_KEY",
    lunarcrush.WithBaseURL("https://lunarcrush.com/api4"), // default, shown for clarity
    lunarcrush.WithTimeout(15*time.Second),
    lunarcrush.WithRetry(3, time.Second),
    lunarcrush.WithHTTPClient(&http.Client{}), // custom transport, proxies, TLS, etc.
)
```

| Option                                                | What it does                                  | When to use it                                      |
|-------------------------------------------------------|-----------------------------------------------|-----------------------------------------------------|
| `WithBaseURL(url string)`                             | Overrides the API base URL.                   | Mocking the API in tests or using a custom gateway. |
| `WithHTTPClient(client *http.Client)`                 | Supplies your own HTTP client.                | Custom TLS, proxies, tracing, or middleware.        |
| `WithTimeout(d time.Duration)`                        | Sets the per-request timeout.                 | Default is 30s; lower it for fast UI endpoints.     |
| `WithRetry(maxAttempts int, baseDelay time.Duration)` | Retries on HTTP 429 with exponential backoff. | Strongly recommended for production workloads.      |

## Exploring the API

The SDK mirrors the LunarCrush API structure. Each resource group is a service you reach through `client`:

```go
client.Coins      // crypto coins
client.Topics     // social topics / hashtags / keywords
client.Stocks     // equities
client.Creators   // social media creators
client.Categories // topic categories
client.Posts      // social posts
client.Searches   // saved searches
client.AI         // AI-generated summaries
client.System     // API changelog / changes
```

Every method accepts `context.Context` first and returns typed response structs plus an `error`. Optional parameters are
pointer fields so you can leave them out without the API receiving unintended zero values.

### Quick example: time-series for a coin

```go
interval := "1w"
bucket := "hour"
series, err := client.Coins.TimeSeries(ctx, "bitcoin", &lunarcrush.TimeSeriesParams{
    Interval: &interval,
    Bucket:   &bucket,
})
if err != nil {
    log.Fatal(err)
}
for _, pt := range series.Data {
    fmt.Printf("time=%d close=%.2f\n", pt.Time, pt.Close)
}
```

### Quick example: creator profile

```go
creator, err := client.Creators.Get(ctx, "twitter", "elonmusk")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s has %.0f followers\n", creator.Data.CreatorDisplayName, creator.Data.FollowerCount)
```

### Quick example: saved search

```go
search, err := client.Searches.Create(ctx, &lunarcrush.SearchCreateParams{
    Name:       "btc-bullish",
    SearchJSON: `{"q":"bitcoin","sentiment":">0.5"}`,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created search: %s\n", search.Data.Slug)
```

## Full API reference

### Topics — `client.Topics`

Understand what the crowd is talking about.

| Method                             | Endpoint                                  | Description                               |
|------------------------------------|-------------------------------------------|-------------------------------------------|
| `Get(ctx, topic)`                  | `GET /public/topic/:topic/v1`             | 24-hour social summary for a topic.       |
| `TimeSeries(ctx, topic, params)`   | `GET /public/topic/:topic/time-series/v1` | Historical social activity over time.     |
| `TimeSeriesV2(ctx, topic, params)` | `GET /public/topic/:topic/time-series/v2` | Same idea, newer format.                  |
| `Creators(ctx, topic)`             | `GET /public/topic/:topic/creators/v1`    | Top creators posting about this topic.    |
| `News(ctx, topic)`                 | `GET /public/topic/:topic/news/v1`        | News articles related to the topic.       |
| `Posts(ctx, topic)`                | `GET /public/topic/:topic/posts/v1`       | Recent social posts for the topic.        |
| `WhatsUp(ctx, topic)`              | `GET /public/topic/:topic/whatsup/v1`     | AI-generated "what is happening" summary. |
| `List(ctx)`                        | `GET /public/topics/list/v1`              | All trending topics.                      |

### Categories — `client.Categories`

Group-related topics under one roof.

| Method                              | Endpoint                                        | Description                           |
|-------------------------------------|-------------------------------------------------|---------------------------------------|
| `List(ctx)`                         | `GET /public/categories/list/v1`                | All available categories.             |
| `Get(ctx, category)`                | `GET /public/category/:category/v1`             | Summary of a single category.         |
| `Creators(ctx, category)`           | `GET /public/category/:category/creators/v1`    | Top creators in the category.         |
| `News(ctx, category)`               | `GET /public/category/:category/news/v1`        | News for the category.                |
| `Posts(ctx, category)`              | `GET /public/category/:category/posts/v1`       | Posts in the category.                |
| `TimeSeries(ctx, category, params)` | `GET /public/category/:category/time-series/v1` | Historical activity for the category. |
| `Topics(ctx, category)`             | `GET /public/category/:category/topics/v1`      | Topics that belong to this category.  |

### Creators — `client.Creators`

Find influencers and their impact.

| Method                                 | Endpoint                                          | Description                    |
|----------------------------------------|---------------------------------------------------|--------------------------------|
| `Get(ctx, network, id)`                | `GET /public/creator/:network/:id/v1`             | Creator profile and metrics.   |
| `Posts(ctx, network, id)`              | `GET /public/creator/:network/:id/posts/v1`       | Recent posts from the creator. |
| `TimeSeries(ctx, network, id, params)` | `GET /public/creator/:network/:id/time-series/v1` | Historical creator activity.   |
| `List(ctx)`                            | `GET /public/creators/list/v1`                    | Global top creators.           |

Supported networks: `twitter`, `youtube`, `instagram`, `reddit`, `tiktok`.

### Posts — `client.Posts`

Raw social content, filtered and aggregated.

| Method                    | Endpoint                           | Description                      |
|---------------------------|------------------------------------|----------------------------------|
| `List(ctx, params)`       | `GET /public/posts/v1`             | List posts by topic or category. |
| `TimeSeries(ctx, params)` | `GET /public/posts/time-series/v1` | Post volume over time.           |

### Coins — `client.Coins`

Crypto market data plus social signals.

| Method                          | Endpoint                                 | Description                           |
|---------------------------------|------------------------------------------|---------------------------------------|
| `List(ctx, params)`             | `GET /public/coins/list/v1`              | List coins (v1).                      |
| `ListV2(ctx, params)`           | `GET /public/coins/list/v2`              | List coins (v2) with extra fields.    |
| `Get(ctx, coin)`                | `GET /public/coins/:coin/v1`             | Single coin details.                  |
| `Meta(ctx, coin)`               | `GET /public/coins/:coin/meta/v1`        | Project metadata, links, description. |
| `TimeSeries(ctx, coin, params)` | `GET /public/coins/:coin/time-series/v2` | OHLCV + social metrics over time.     |

### Stocks — `client.Stocks`

Traditional equities through the LunarCrush lens.

| Method                           | Endpoint                                   | Description                     |
|----------------------------------|--------------------------------------------|---------------------------------|
| `List(ctx, params)`              | `GET /public/stocks/list/v1`               | List stocks (v1).               |
| `ListV2(ctx, params)`            | `GET /public/stocks/list/v2`               | List stocks (v2).               |
| `Get(ctx, stock)`                | `GET /public/stocks/:stock/v1`             | Single stock details.           |
| `TimeSeries(ctx, stock, params)` | `GET /public/stocks/:stock/time-series/v2` | Historical stock + social data. |

### Searches — `client.Searches`

Save and run custom queries.

| Method                      | Endpoint                            | Description                          |
|-----------------------------|-------------------------------------|--------------------------------------|
| `Create(ctx, params)`       | `GET /public/searches/create`       | Create a new saved search.           |
| `List(ctx)`                 | `GET /public/searches/list`         | List all saved searches.             |
| `Search(ctx, searchJSON)`   | `GET /public/searches/search`       | Run an ad-hoc search without saving. |
| `Get(ctx, slug)`            | `GET /public/searches/:slug`        | Fetch a saved search.                |
| `Update(ctx, slug, params)` | `GET /public/searches/:slug/update` | Update a saved search.               |
| `Delete(ctx, slug)`         | `GET /public/searches/:slug/delete` | Delete a saved search.               |

### AI — `client.AI`

Natural-language summaries from LunarCrush AI.

| Method                      | Endpoint                              | Description              |
|-----------------------------|---------------------------------------|--------------------------|
| `Topic(ctx, topic)`         | `GET /public/ai/topic/:topic`         | AI summary of a topic.   |
| `Creator(ctx, network, id)` | `GET /public/ai/creator/:network/:id` | AI summary of a creator. |

### System — `client.System`

Keep up with API changes.

| Method         | Endpoint                     | Description                               |
|----------------|------------------------------|-------------------------------------------|
| `Changes(ctx)` | `GET /public/system/changes` | Recent LunarCrush API / platform changes. |

## Rate limits and retry behavior

LunarCrush rate limits depend on your plan:

| Plan       | Requests/min | Requests/day |
|------------|--------------|--------------|
| Hobby      | 4            | 100          |
| Individual | 10           | 2,000        |
| Builder    | 100          | 20,000       |
| Scale      | 500          | 100,000      |

Hitting a 429 is not a panic moment. Enable retries and the SDK will back off automatically, doubling the wait time on
each attempt. If the API sends a `Retry-After` header, we honor that instead of guessing.

```go
client := lunarcrush.NewClient("YOUR_API_KEY",
    lunarcrush.WithRetry(3, time.Second), // 1s, then 2s, then 4s
)
```

**Tip:** If you are on the Hobby plan, a 1-minute backoff is often too aggressive; use `WithRetry` to keep your program
polite and your API quota healthy.

## Error handling done right

Every non-2xx response comes back as an `*lunarcrush.APIError`:

```go
type APIError struct {
    StatusCode int
    Message    string
    RawBody    []byte
}
```

For the most common HTTP statuses, you can use `errors.Is` with sentinel errors:

```go
import "errors"

topic, err := client.Topics.Get(ctx, "totally-unknown-topic")
if err != nil {
    switch {
    case errors.Is(err, lunarcrush.ErrUnauthorized):
        log.Println("Check your API key — LunarCrush rejected it.")
    case errors.Is(err, lunarcrush.ErrNotFound):
        log.Println("That topic or resource does not exist.")
    case errors.Is(err, lunarcrush.ErrRateLimited):
        log.Println("Rate limit hit and retries exhausted. Slow down or upgrade your plan.")
    default:
        var apiErr *lunarcrush.APIError
        if errors.As(err, &apiErr) {
            log.Printf("LunarCrush returned HTTP %d: %s", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

Because we keep the raw response body in `APIError.RawBody`, debugging weird responses is much easier than a generic
error string.

## Contexts and concurrency

All methods are context-first. Pass a `context.Background()`, a `context.WithTimeout()`, or a `context.WithCancel()` —
whatever fits your flow:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

coin, err := client.Coins.Get(ctx, "bitcoin")
```

The client is also safe to share across goroutines. Fetch several coins in parallel without creating new clients:

```go
var wg sync.WaitGroup
for _, sym := range []string{"bitcoin", "ethereum", "solana", "cardano"} {
    wg.Add(1)
    go func(symbol string) {
        defer wg.Done()
        coin, err := client.Coins.Get(ctx, symbol)
        if err != nil {
            log.Printf("%s failed: %v", symbol, err)
            return
        }
        log.Printf("%s galaxy score: %.1f", symbol, coin.Data.GalaxyScore)
    }(sym)
}
wg.Wait()
```

## A real-world example: nightly market report

Imagine you want to email a short daily report every morning. Here is the skeleton:

```go
ctx := context.Background()

sort := "galaxy_score"
limit := 5
coins, err := client.Coins.List(ctx, &lunarcrush.CoinsListParams{
    Sort:  &sort,
    Limit: &limit,
    Desc:  ptr(true),
})
if err != nil {
    log.Fatal(err)
}

fmt.Println("Top 5 coins by Galaxy Score")
for _, c := range coins.Data {
    fmt.Printf("- %s (%s): score %.1f, price $%.2f, 24h change %.2f%%\n",
        c.Name, c.Symbol, c.GalaxyScore, c.Price, c.PercentChange24h)
}

stock, err := client.Stocks.Get(ctx, "NVDA")
if err == nil {
    fmt.Printf("\nNVDA social activity: %.0f interactions (24h)\n", stock.Data.Interactions24h)
}
```

## Best practices

1. **Reuse the client.** Creating one `*Client` and sharing it is cheaper and cleaner than making a new one per request.
2. **Always pass a context.** Timeouts and cancellation make your application robust.
3. **Enable retries in production.** LunarCrush rate limits are real; `WithRetry` handles them gracefully.
4. **Check sentinel errors first.** They are the fastest way to react to 401, 404, and 429.
5. **Use `limit` and `sort`.** Many list endpoints support them; they keep responses small and focused.
6. **Pointer fields for optional params.** This is how the SDK distinguishes "not set" from "zero value".

## Running the tests

The SDK is fully tested with `net/http/httptest`. Clone it and run:

```bash
go test ./...
```

You will see coverage for JSON decoding, retry logic, context cancellation, error mapping, and every service endpoint.

## Troubleshooting

| Symptom               | Likely cause                                     | Fix                                                                      |
|-----------------------|--------------------------------------------------|--------------------------------------------------------------------------|
| `401 Unauthorized`    | Invalid or missing API key.                      | Verify `YOUR_API_KEY` in `NewClient`.                                    |
| `404 Not Found`       | The topic/coin/stock/creator does not exist.     | Check the exact identifier LunarCrush expects.                           |
| `429 Rate Limited`    | You hit the plan limit.                          | Add `WithRetry`, reduce request frequency, or upgrade your plan.         |
| Timeouts              | Request took longer than the configured timeout. | Increase `WithTimeout`, or use a `context.WithTimeout`.                  |
| Empty response fields | LunarCrush omitted them for this request.        | Fields are tagged with `omitempty`; not every call returns every metric. |

## Contributing

Found a bug? Have an idea for a better example? Pull requests are welcome!

1. Fork the repo at [github.com/tigusigalpa/lunarcrush-go](https://github.com/tigusigalpa/lunarcrush-go).
2. Create a feature branch.
3. Add or update tests for any changed behavior.
4. Run `go vet ./...` and `go test ./...` before pushing.

## License

MIT © [Igor Sazonov](https://github.com/tigusigalpa)

Happy building! 🚀
