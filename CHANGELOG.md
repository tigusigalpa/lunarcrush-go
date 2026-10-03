# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Changed

- **Breaking:** nullable `Coin.MaxSupply`, post/news fields, and time-series
  post counts now use pointers so API `null` values are not silently coerced to
  zero values. Update consumers to handle `nil`.
- Added raw response receipts to selected coin and topic response envelopes for
  provenance, unknown fields, and lossless numeric inspection via
  `RawResponse`.
- Added missing topic fields, including related topics, sentiment detail, prior
  rank horizons, and post counts.

## [1.0.0] - 2025-01-01

### Added

- Initial release of the LunarCrush Go SDK.
- `Client` with functional options: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithRetry`.
- `CoinsService` — list coins (v1/v2), get coin, coin meta, coin time-series.
- `TopicsService` — topic summary, time-series (v1/v2), creators, news, posts, whatsup, list.
- `StocksService` — list stocks (v1/v2), get stock, stock time-series.
- `CreatorsService` — get creator, creator posts, creator time-series, list creators.
- `CategoriesService` — list categories, get category, category creators/news/posts/time-series/topics.
- `PostsService` — list posts, posts time-series.
- `SearchesService` — create, list, search, get, update, delete searches.
- `SystemService` — system changes.
- `AIService` — AI topic and AI creator summaries.
- Automatic retry with exponential backoff for HTTP 429 responses, honoring `Retry-After`.
- Sentinel errors: `ErrUnauthorized`, `ErrNotFound`, `ErrRateLimited`.
- Full test suite using `net/http/httptest`.
