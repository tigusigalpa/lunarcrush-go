package lunarcrush

import (
	"context"
	"fmt"
)

// CreatorsService provides access to the /public/creator* and
// /public/creators* endpoints of the LunarCrush API.
type CreatorsService struct {
	client *Client
}

// CreatorProfile represents a detailed social media creator profile as
// returned by GET /public/creator/:network/:id/v1.
type CreatorProfile struct {
	CreatorID          string  `json:"creator_id,omitempty"`
	CreatorName        string  `json:"creator_name,omitempty"`
	CreatorDisplayName string  `json:"creator_display_name,omitempty"`
	CreatorAvatar      string  `json:"creator_avatar,omitempty"`
	Network            string  `json:"creator_network,omitempty"`
	FollowerCount      float64 `json:"creator_followers,omitempty"`
	Interactions24h    float64 `json:"interactions_24h,omitempty"`
	PostCount24h       int     `json:"creator_posts_24h,omitempty"`
	Rank               int     `json:"creator_rank,omitempty"`
	TopicsCovered      []string `json:"topics,omitempty"`
}

// CreatorResponse is the response envelope returned when fetching a
// single creator profile.
type CreatorResponse struct {
	Data   CreatorProfile         `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CreatorPostsResponse is the response envelope returned by the creator
// posts endpoint.
type CreatorPostsResponse struct {
	Data   []Post                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CreatorTimeSeriesResponse is the response envelope returned by the
// creator time-series endpoint.
type CreatorTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// CreatorsListResponse is the response envelope returned by
// GET /public/creators/list/v1.
type CreatorsListResponse struct {
	Data   []Creator              `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Get retrieves the profile for a single creator on a given social
// network: GET /public/creator/:network/:id/v1. Supported networks
// include "twitter", "youtube", "instagram", "reddit", and "tiktok".
func (s *CreatorsService) Get(ctx context.Context, network, id string) (*CreatorResponse, error) {
	var out CreatorResponse
	path := fmt.Sprintf("/public/creator/%s/%s/v1", network, id)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Posts retrieves recent posts authored by a creator:
// GET /public/creator/:network/:id/posts/v1.
func (s *CreatorsService) Posts(ctx context.Context, network, id string) (*CreatorPostsResponse, error) {
	var out CreatorPostsResponse
	path := fmt.Sprintf("/public/creator/%s/%s/posts/v1", network, id)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for a creator:
// GET /public/creator/:network/:id/time-series/v1.
func (s *CreatorsService) TimeSeries(ctx context.Context, network, id string, params *TimeSeriesParams) (*CreatorTimeSeriesResponse, error) {
	var out CreatorTimeSeriesResponse
	path := fmt.Sprintf("/public/creator/%s/%s/time-series/v1", network, id)
	if err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List retrieves the list of top creators across all networks:
// GET /public/creators/list/v1.
func (s *CreatorsService) List(ctx context.Context) (*CreatorsListResponse, error) {
	var out CreatorsListResponse
	if err := s.client.doRequest(ctx, "GET", "/public/creators/list/v1", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
