package lunarcrush

import (
	"context"
)

// PostsService provides access to the /public/posts* endpoints of the
// LunarCrush API.
type PostsService struct {
	client *Client
}

// PostsListParams holds optional query parameters accepted by the posts
// list endpoint.
type PostsListParams struct {
	// Topic restricts results to posts related to a given topic.
	Topic *string
	// Category restricts results to posts related to a given category.
	Category *string
	// Limit restricts the number of returned results.
	Limit *int
}

// toQuery converts PostsListParams into a queryValues helper.
func (p *PostsListParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("topic", p.Topic)
	q.setString("category", p.Category)
	q.setInt("limit", p.Limit)
	return q
}

// PostsListResponse is the response envelope returned by
// GET /public/posts/v1.
type PostsListResponse struct {
	Data   []Post                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// PostsTimeSeriesParams holds optional query parameters accepted by the
// posts time-series endpoint.
type PostsTimeSeriesParams struct {
	// Topic restricts results to posts related to a given topic.
	Topic *string
	// Category restricts results to posts related to a given category.
	Category *string
	// Bucket selects the aggregation bucket, e.g. "hour" or "day".
	Bucket *string
	// Interval selects a predefined time window, e.g. "1w", "1m", "3m".
	Interval *string
	// Start is a Unix timestamp specifying the beginning of the range.
	Start *int
	// End is a Unix timestamp specifying the end of the range.
	End *int
}

// toQuery converts PostsTimeSeriesParams into a queryValues helper.
func (p *PostsTimeSeriesParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("topic", p.Topic)
	q.setString("category", p.Category)
	q.setString("bucket", p.Bucket)
	q.setString("interval", p.Interval)
	q.setInt("start", p.Start)
	q.setInt("end", p.End)
	return q
}

// PostsTimeSeriesResponse is the response envelope returned by
// GET /public/posts/time-series/v1.
type PostsTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// List retrieves a list of social posts: GET /public/posts/v1.
func (s *PostsService) List(ctx context.Context, params *PostsListParams) (*PostsListResponse, error) {
	var out PostsListResponse
	err := s.client.doRequest(ctx, "GET", "/public/posts/v1", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for posts activity:
// GET /public/posts/time-series/v1.
func (s *PostsService) TimeSeries(ctx context.Context, params *PostsTimeSeriesParams) (*PostsTimeSeriesResponse, error) {
	var out PostsTimeSeriesResponse
	err := s.client.doRequest(ctx, "GET", "/public/posts/time-series/v1", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
