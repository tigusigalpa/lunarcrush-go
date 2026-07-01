package lunarcrush

import (
	"context"
	"fmt"
)

// CategoriesService provides access to the /public/categories* and
// /public/category* endpoints of the LunarCrush API.
type CategoriesService struct {
	client *Client
}

// Category represents a topic category summary as returned by the
// LunarCrush API.
type Category struct {
	Category        string  `json:"category,omitempty"`
	Title           string  `json:"title,omitempty"`
	CategoryRank    int     `json:"category_rank,omitempty"`
	Interactions24h float64 `json:"interactions_24h,omitempty"`
	NumContributors int     `json:"num_contributors,omitempty"`
	NumPosts        int     `json:"num_posts,omitempty"`
	SocialVolume24h float64 `json:"social_volume_24h,omitempty"`
	SocialDominance float64 `json:"social_dominance,omitempty"`
	Sentiment       float64 `json:"sentiment,omitempty"`
	TrendDirection  string  `json:"trend,omitempty"`
}

// CategoriesListResponse is the response envelope returned by
// GET /public/categories/list/v1.
type CategoriesListResponse struct {
	Data   []Category             `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryResponse is the response envelope returned when fetching a
// single category.
type CategoryResponse struct {
	Data   Category               `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryCreatorsResponse is the response envelope returned by the
// category creators endpoint.
type CategoryCreatorsResponse struct {
	Data   []Creator              `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryNewsResponse is the response envelope returned by the category
// news endpoint.
type CategoryNewsResponse struct {
	Data   []NewsItem             `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryPostsResponse is the response envelope returned by the
// category posts endpoint.
type CategoryPostsResponse struct {
	Data   []Post                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryTimeSeriesResponse is the response envelope returned by the
// category time-series endpoint.
type CategoryTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CategoryTopicsResponse is the response envelope returned by the
// category topics endpoint.
type CategoryTopicsResponse struct {
	Data   []TopicListItem        `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// List retrieves the list of all categories:
// GET /public/categories/list/v1.
func (s *CategoriesService) List(ctx context.Context) (*CategoriesListResponse, error) {
	var out CategoriesListResponse
	if err := s.client.doRequest(ctx, "GET", "/public/categories/list/v1", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves the summary for a single category:
// GET /public/category/:category/v1.
func (s *CategoriesService) Get(ctx context.Context, category string) (*CategoryResponse, error) {
	var out CategoryResponse
	path := fmt.Sprintf("/public/category/%s/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Creators retrieves the top creators posting about a category:
// GET /public/category/:category/creators/v1.
func (s *CategoriesService) Creators(ctx context.Context, category string) (*CategoryCreatorsResponse, error) {
	var out CategoryCreatorsResponse
	path := fmt.Sprintf("/public/category/%s/creators/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// News retrieves news articles related to a category:
// GET /public/category/:category/news/v1.
func (s *CategoriesService) News(ctx context.Context, category string) (*CategoryNewsResponse, error) {
	var out CategoryNewsResponse
	path := fmt.Sprintf("/public/category/%s/news/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Posts retrieves social posts related to a category:
// GET /public/category/:category/posts/v1.
func (s *CategoriesService) Posts(ctx context.Context, category string) (*CategoryPostsResponse, error) {
	var out CategoryPostsResponse
	path := fmt.Sprintf("/public/category/%s/posts/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for a category:
// GET /public/category/:category/time-series/v1.
func (s *CategoriesService) TimeSeries(ctx context.Context, category string, params *TimeSeriesParams) (*CategoryTimeSeriesResponse, error) {
	var out CategoryTimeSeriesResponse
	path := fmt.Sprintf("/public/category/%s/time-series/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Topics retrieves the topics that make up a category:
// GET /public/category/:category/topics/v1.
func (s *CategoriesService) Topics(ctx context.Context, category string) (*CategoryTopicsResponse, error) {
	var out CategoryTopicsResponse
	path := fmt.Sprintf("/public/category/%s/topics/v1", category)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
