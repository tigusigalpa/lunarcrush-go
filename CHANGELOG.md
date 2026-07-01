# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
