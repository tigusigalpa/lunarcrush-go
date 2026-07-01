package lunarcrush

import (
	"context"
	"fmt"
)

// CoinsService provides access to the /public/coins* endpoints of the
// LunarCrush API, covering coin listings and individual coin data.
type CoinsService struct {
	client *Client
}

// Coin represents a single cryptocurrency entry as returned by the
// LunarCrush API, including market and social metrics.
type Coin struct {
	ID                   int     `json:"id,omitempty"`
	Symbol               string  `json:"symbol,omitempty"`
	Name                 string  `json:"name,omitempty"`
	Price                float64 `json:"price,omitempty"`
	PriceBTC             float64 `json:"price_btc,omitempty"`
	Volume24h            float64 `json:"volume_24h,omitempty"`
	VolatilityDay        float64 `json:"volatility,omitempty"`
	CirculatingSupply    float64 `json:"circulating_supply,omitempty"`
	MaxSupply            float64 `json:"max_supply,omitempty"`
	PercentChange1h      float64 `json:"percent_change_1h,omitempty"`
	PercentChange24h     float64 `json:"percent_change_24h,omitempty"`
	PercentChange7d      float64 `json:"percent_change_7d,omitempty"`
	PercentChange30d     float64 `json:"percent_change_30d,omitempty"`
	MarketCap            float64 `json:"market_cap,omitempty"`
	MarketCapRank        int     `json:"market_cap_rank,omitempty"`
	Interactions24h      float64 `json:"interactions_24h,omitempty"`
	SocialVolume24h      float64 `json:"social_volume_24h,omitempty"`
	SocialDominance      float64 `json:"social_dominance,omitempty"`
	MarketDominance      float64 `json:"market_dominance,omitempty"`
	GalaxyScore          float64 `json:"galaxy_score,omitempty"`
	AltRank              int     `json:"alt_rank,omitempty"`
	SentimentRelative    float64 `json:"sentiment,omitempty"`
	Categories           string  `json:"categories,omitempty"`
	Blockchains          []struct {
		Type    string `json:"type,omitempty"`
		Network string `json:"network,omitempty"`
		Address string `json:"address,omitempty"`
		Decimals int   `json:"decimals,omitempty"`
	} `json:"blockchains,omitempty"`
	LastUpdatedPrice float64 `json:"last_updated_price,omitempty"`
	Logo             string  `json:"logo,omitempty"`
}

// CoinListResponse is the response envelope returned by the coins list
// endpoints (v1 and v2).
type CoinListResponse struct {
	Data   []Coin                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CoinResponse is the response envelope returned when fetching a single
// coin.
type CoinResponse struct {
	Data   Coin                   `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CoinMeta contains descriptive metadata about a coin such as its
// website links, description, and social handles.
type CoinMeta struct {
	ID           int      `json:"id,omitempty"`
	Symbol       string   `json:"symbol,omitempty"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Github       []string `json:"github_link,omitempty"`
	Website      []string `json:"website_link,omitempty"`
	Whitepaper   string   `json:"whitepaper_link,omitempty"`
	Twitter      string   `json:"twitter_link,omitempty"`
	Reddit       string   `json:"reddit_link,omitempty"`
	Blog         []string `json:"blog_link,omitempty"`
	Facebook     string   `json:"facebook_link,omitempty"`
}

// CoinMetaResponse is the response envelope returned by the coin meta
// endpoint.
type CoinMetaResponse struct {
	Data   CoinMeta               `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// TimeSeriesPoint represents a single bucketed data point in a
// time-series response, containing a Unix timestamp plus a flexible set
// of numeric metrics.
type TimeSeriesPoint struct {
	Time             int64   `json:"time,omitempty"`
	Open             float64 `json:"open,omitempty"`
	Close            float64 `json:"close,omitempty"`
	High             float64 `json:"high,omitempty"`
	Low              float64 `json:"low,omitempty"`
	Volume24h        float64 `json:"volume_24h,omitempty"`
	MarketCap        float64 `json:"market_cap,omitempty"`
	Interactions     float64 `json:"interactions,omitempty"`
	SocialVolume     float64 `json:"social_volume,omitempty"`
	SocialDominance  float64 `json:"social_dominance,omitempty"`
	GalaxyScore      float64 `json:"galaxy_score,omitempty"`
	AltRank          int     `json:"alt_rank,omitempty"`
	Sentiment        float64 `json:"sentiment,omitempty"`
	SpamCount        float64 `json:"spam,omitempty"`
	Posts            float64 `json:"posts_created,omitempty"`
	Contributors     float64 `json:"contributors_active,omitempty"`
}

// CoinTimeSeriesResponse is the response envelope returned by the coin
// time-series endpoint.
type CoinTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CoinsListParams holds optional query parameters accepted by the coins
// list endpoints. All fields are optional; nil/zero fields are omitted
// from the request.
type CoinsListParams struct {
	// Sort specifies the field used to sort results, e.g. "galaxy_score".
	Sort *string
	// Desc, when true, sorts results in descending order.
	Desc *bool
	// Limit restricts the number of returned results.
	Limit *int
	// Page selects a page of results when the endpoint supports pagination.
	Page *int
	// Filter restricts results to coins matching a category/filter value.
	Filter *string
}

// toQuery converts CoinsListParams into a url.Values-compatible helper.
func (p *CoinsListParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("sort", p.Sort)
	q.setBool("desc", p.Desc)
	q.setInt("limit", p.Limit)
	q.setInt("page", p.Page)
	q.setString("filter", p.Filter)
	return q
}

// List retrieves a list of coins ranked and sorted according to params
// using the v1 endpoint: GET /public/coins/list/v1.
func (s *CoinsService) List(ctx context.Context, params *CoinsListParams) (*CoinListResponse, error) {
	var out CoinListResponse
	err := s.client.doRequest(ctx, "GET", "/public/coins/list/v1", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ListV2 retrieves a list of coins using the v2 endpoint, which typically
// includes additional fields compared to v1: GET /public/coins/list/v2.
func (s *CoinsService) ListV2(ctx context.Context, params *CoinsListParams) (*CoinListResponse, error) {
	var out CoinListResponse
	err := s.client.doRequest(ctx, "GET", "/public/coins/list/v2", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves details for a single coin by its symbol or id:
// GET /public/coins/:coin/v1.
func (s *CoinsService) Get(ctx context.Context, coin string) (*CoinResponse, error) {
	var out CoinResponse
	path := fmt.Sprintf("/public/coins/%s/v1", coin)
	err := s.client.doRequest(ctx, "GET", path, nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Meta retrieves descriptive metadata for a coin:
// GET /public/coins/:coin/meta/v1.
func (s *CoinsService) Meta(ctx context.Context, coin string) (*CoinMetaResponse, error) {
	var out CoinMetaResponse
	path := fmt.Sprintf("/public/coins/%s/meta/v1", coin)
	err := s.client.doRequest(ctx, "GET", path, nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for a coin:
// GET /public/coins/:coin/time-series/v2.
func (s *CoinsService) TimeSeries(ctx context.Context, coin string, params *TimeSeriesParams) (*CoinTimeSeriesResponse, error) {
	var out CoinTimeSeriesResponse
	path := fmt.Sprintf("/public/coins/%s/time-series/v2", coin)
	err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
