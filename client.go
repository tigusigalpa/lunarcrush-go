// Package lunarcrush provides a Go SDK for the LunarCrush API v4
// (https://lunarcrush.com/api4). It exposes service structs grouped by
// resource (Coins, Topics, Stocks, Creators, Categories, Posts, Searches,
// System, AI) accessible as fields on Client.
package lunarcrush

import (
	"net/http"
	"time"
)

// DefaultBaseURL is the default LunarCrush API v4 base URL used by
// NewClient unless overridden with WithBaseURL.
const DefaultBaseURL = "https://lunarcrush.com/api4"

// DefaultTimeout is the default HTTP client timeout used unless
// overridden with WithTimeout.
const DefaultTimeout = 30 * time.Second

// DefaultMaxAttempts is the default number of attempts (including the
// initial request) made for a request before giving up when receiving
// HTTP 429 responses, unless overridden with WithRetry.
const DefaultMaxAttempts = 1

// DefaultBaseDelay is the default initial backoff delay used for retry
// attempts, unless overridden with WithRetry.
const DefaultBaseDelay = 500 * time.Millisecond

// Client is the entry point of the LunarCrush Go SDK. It holds the HTTP
// configuration and exposes one service struct per API resource group.
// A Client is safe for concurrent use by multiple goroutines.
type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	maxAttempts int
	baseDelay   time.Duration

	// Coins provides access to the /public/coins* endpoints.
	Coins *CoinsService
	// Topics provides access to the /public/topic* and /public/topics* endpoints.
	Topics *TopicsService
	// Stocks provides access to the /public/stocks* endpoints.
	Stocks *StocksService
	// Creators provides access to the /public/creator* and /public/creators* endpoints.
	Creators *CreatorsService
	// Categories provides access to the /public/category* and /public/categories* endpoints.
	Categories *CategoriesService
	// Posts provides access to the /public/posts* endpoints.
	Posts *PostsService
	// Searches provides access to the /public/searches* endpoints.
	Searches *SearchesService
	// System provides access to the /public/system* endpoints.
	System *SystemService
	// AI provides access to the /public/ai* endpoints.
	AI *AIService
}

// NewClient creates a new LunarCrush API client authenticated with the
// given apiKey. The apiKey is sent as a Bearer token in the Authorization
// header on every request. Behaviour can be customized via Option values
// such as WithBaseURL, WithHTTPClient, WithTimeout, and WithRetry.
//
// Example:
//
//	client := lunarcrush.NewClient("YOUR_API_KEY",
//	    lunarcrush.WithTimeout(15*time.Second),
//	    lunarcrush.WithRetry(3, time.Second),
//	)
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:      apiKey,
		baseURL:     DefaultBaseURL,
		httpClient:  &http.Client{Timeout: DefaultTimeout},
		maxAttempts: DefaultMaxAttempts,
		baseDelay:   DefaultBaseDelay,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Coins = &CoinsService{client: c}
	c.Topics = &TopicsService{client: c}
	c.Stocks = &StocksService{client: c}
	c.Creators = &CreatorsService{client: c}
	c.Categories = &CategoriesService{client: c}
	c.Posts = &PostsService{client: c}
	c.Searches = &SearchesService{client: c}
	c.System = &SystemService{client: c}
	c.AI = &AIService{client: c}

	return c
}
