package lunarcrush

import (
	"context"
	"fmt"
)

// StocksService provides access to the /public/stocks* endpoints of the
// LunarCrush API.
type StocksService struct {
	client *Client
}

// Stock represents a single stock/equity entry as returned by the
// LunarCrush API, including market and social metrics.
type Stock struct {
	ID               int     `json:"id,omitempty"`
	Symbol           string  `json:"symbol,omitempty"`
	Name             string  `json:"name,omitempty"`
	Price            float64 `json:"price,omitempty"`
	Volume24h        float64 `json:"volume_24h,omitempty"`
	PercentChange24h float64 `json:"percent_change_24h,omitempty"`
	MarketCap        float64 `json:"market_cap,omitempty"`
	MarketCapRank    int     `json:"market_cap_rank,omitempty"`
	Interactions24h  float64 `json:"interactions_24h,omitempty"`
	SocialVolume24h  float64 `json:"social_volume_24h,omitempty"`
	SocialDominance  float64 `json:"social_dominance,omitempty"`
	GalaxyScore      float64 `json:"galaxy_score,omitempty"`
	AltRank          int     `json:"alt_rank,omitempty"`
	Sentiment        float64 `json:"sentiment,omitempty"`
}

// StockListResponse is the response envelope returned by the stocks
// list endpoints (v1 and v2).
type StockListResponse struct {
	Data   []Stock                `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// StockResponse is the response envelope returned when fetching a single
// stock.
type StockResponse struct {
	Data   Stock                  `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// StockTimeSeriesResponse is the response envelope returned by the stock
// time-series endpoint.
type StockTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// StocksListParams holds optional query parameters accepted by the
// stocks list endpoints.
type StocksListParams struct {
	// Sort specifies the field used to sort results.
	Sort *string
	// Desc, when true, sorts results in descending order.
	Desc *bool
	// Limit restricts the number of returned results.
	Limit *int
	// Page selects a page of results when the endpoint supports pagination.
	Page *int
}

// toQuery converts StocksListParams into a queryValues helper.
func (p *StocksListParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("sort", p.Sort)
	q.setBool("desc", p.Desc)
	q.setInt("limit", p.Limit)
	q.setInt("page", p.Page)
	return q
}

// List retrieves a list of stocks using the v1 endpoint:
// GET /public/stocks/list/v1.
func (s *StocksService) List(ctx context.Context, params *StocksListParams) (*StockListResponse, error) {
	var out StockListResponse
	err := s.client.doRequest(ctx, "GET", "/public/stocks/list/v1", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ListV2 retrieves a list of stocks using the v2 endpoint:
// GET /public/stocks/list/v2.
func (s *StocksService) ListV2(ctx context.Context, params *StocksListParams) (*StockListResponse, error) {
	var out StockListResponse
	err := s.client.doRequest(ctx, "GET", "/public/stocks/list/v2", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves details for a single stock by its ticker symbol:
// GET /public/stocks/:stock/v1.
func (s *StocksService) Get(ctx context.Context, stock string) (*StockResponse, error) {
	var out StockResponse
	path := fmt.Sprintf("/public/stocks/%s/v1", stock)
	err := s.client.doRequest(ctx, "GET", path, nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for a stock:
// GET /public/stocks/:stock/time-series/v2.
func (s *StocksService) TimeSeries(ctx context.Context, stock string, params *TimeSeriesParams) (*StockTimeSeriesResponse, error) {
	var out StockTimeSeriesResponse
	path := fmt.Sprintf("/public/stocks/%s/time-series/v2", stock)
	err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
