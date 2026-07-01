package lunarcrush

import (
	"context"
	"fmt"
)

// TopicsService provides access to the /public/topic* and
// /public/topics* endpoints of the LunarCrush API.
type TopicsService struct {
	client *Client
}

// TimeSeriesParams holds optional query parameters shared by every
// time-series endpoint across services (Coins, Topics, Stocks, Creators,
// Categories). All fields are optional.
type TimeSeriesParams struct {
	// Bucket selects the aggregation bucket, e.g. "hour" or "day".
	Bucket *string
	// Interval selects a predefined time window, e.g. "1w", "1m", "3m".
	Interval *string
	// Start is a Unix timestamp specifying the beginning of the range.
	Start *int
	// End is a Unix timestamp specifying the end of the range.
	End *int
}

// toQuery converts TimeSeriesParams into a queryValues helper.
func (p *TimeSeriesParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("bucket", p.Bucket)
	q.setString("interval", p.Interval)
	q.setInt("start", p.Start)
	q.setInt("end", p.End)
	return q
}

// Topic represents the 24-hour social activity summary for a topic
// (e.g. a keyword, hashtag, or ticker symbol).
type Topic struct {
	Topic              string  `json:"topic,omitempty"`
	Title              string  `json:"title,omitempty"`
	TopicRank          int     `json:"topic_rank,omitempty"`
	Interactions24h    float64 `json:"interactions_24h,omitempty"`
	NumContributors    int     `json:"num_contributors,omitempty"`
	NumPosts           int     `json:"num_posts,omitempty"`
	SocialVolume24h    float64 `json:"social_volume_24h,omitempty"`
	SocialDominance    float64 `json:"social_dominance,omitempty"`
	Sentiment          float64 `json:"sentiment,omitempty"`
	Categories         []string `json:"categories,omitempty"`
	TrendDirection     string  `json:"trend,omitempty"`
	TypesCount         map[string]int     `json:"types_count,omitempty"`
	TypesInteractions  map[string]float64 `json:"types_interactions,omitempty"`
	TypesSentiment     map[string]float64 `json:"types_sentiment,omitempty"`
}

// TopicResponse is the response envelope returned by the topic summary
// endpoint.
type TopicResponse struct {
	Data   Topic                  `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// TopicTimeSeriesResponse is the response envelope returned by the topic
// time-series endpoints (v1 and v2).
type TopicTimeSeriesResponse struct {
	Data   []TimeSeriesPoint      `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Creator represents a social media creator/influencer profile as
// returned by topic- and creator-related endpoints.
type Creator struct {
	CreatorID          string  `json:"creator_id,omitempty"`
	CreatorName        string  `json:"creator_name,omitempty"`
	CreatorDisplayName string  `json:"creator_display_name,omitempty"`
	CreatorAvatar      string  `json:"creator_avatar,omitempty"`
	Network            string  `json:"creator_network,omitempty"`
	FollowerCount      float64 `json:"creator_followers,omitempty"`
	Interactions24h    float64 `json:"interactions_24h,omitempty"`
	PostCount          int     `json:"creator_posts,omitempty"`
	Rank               int     `json:"creator_rank,omitempty"`
}

// TopicCreatorsResponse is the response envelope returned by the topic
// creators endpoint.
type TopicCreatorsResponse struct {
	Data   []Creator              `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// NewsItem represents a news article associated with a topic or
// category.
type NewsItem struct {
	ID              string  `json:"id,omitempty"`
	PostTitle       string  `json:"post_title,omitempty"`
	PostLink        string  `json:"post_link,omitempty"`
	PostType        string  `json:"post_type,omitempty"`
	PostSentiment   float64 `json:"post_sentiment,omitempty"`
	Interactions24h float64 `json:"interactions_24h,omitempty"`
	CreatedTime     int64   `json:"post_created,omitempty"`
	CreatorName     string  `json:"creator_name,omitempty"`
}

// TopicNewsResponse is the response envelope returned by the topic news
// endpoint.
type TopicNewsResponse struct {
	Data   []NewsItem             `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Post represents a social media post associated with a topic, category,
// or creator.
type Post struct {
	ID              string  `json:"id,omitempty"`
	PostType        string  `json:"post_type,omitempty"`
	PostTitle       string  `json:"post_title,omitempty"`
	PostLink        string  `json:"post_link,omitempty"`
	PostSentiment   float64 `json:"post_sentiment,omitempty"`
	Interactions24h float64 `json:"interactions_24h,omitempty"`
	CreatedTime     int64   `json:"post_created,omitempty"`
	CreatorID       string  `json:"creator_id,omitempty"`
	CreatorName     string  `json:"creator_name,omitempty"`
	CreatorNetwork  string  `json:"creator_network,omitempty"`
}

// TopicPostsResponse is the response envelope returned by the topic
// posts endpoint.
type TopicPostsResponse struct {
	Data   []Post                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// WhatsUp represents an AI-generated natural-language summary of recent
// activity for a topic.
type WhatsUp struct {
	Topic   string `json:"topic,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// TopicWhatsUpResponse is the response envelope returned by the topic
// whatsup endpoint.
type TopicWhatsUpResponse struct {
	Data   WhatsUp                `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// TopicListItem represents a summarized entry from the topics list
// endpoint.
type TopicListItem struct {
	Topic           string  `json:"topic,omitempty"`
	Title           string  `json:"title,omitempty"`
	TopicRank       int     `json:"topic_rank,omitempty"`
	Interactions24h float64 `json:"interactions_24h,omitempty"`
	NumContributors int     `json:"num_contributors,omitempty"`
}

// TopicsListResponse is the response envelope returned by
// GET /public/topics/list/v1.
type TopicsListResponse struct {
	Data   []TopicListItem        `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Get retrieves the 24-hour social summary for a topic:
// GET /public/topic/:topic/v1.
func (s *TopicsService) Get(ctx context.Context, topic string) (*TopicResponse, error) {
	var out TopicResponse
	path := fmt.Sprintf("/public/topic/%s/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeries retrieves historical time-series data for a topic using the
// v1 endpoint: GET /public/topic/:topic/time-series/v1.
func (s *TopicsService) TimeSeries(ctx context.Context, topic string, params *TimeSeriesParams) (*TopicTimeSeriesResponse, error) {
	var out TopicTimeSeriesResponse
	path := fmt.Sprintf("/public/topic/%s/time-series/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TimeSeriesV2 retrieves historical time-series data for a topic using
// the v2 endpoint: GET /public/topic/:topic/time-series/v2.
func (s *TopicsService) TimeSeriesV2(ctx context.Context, topic string, params *TimeSeriesParams) (*TopicTimeSeriesResponse, error) {
	var out TopicTimeSeriesResponse
	path := fmt.Sprintf("/public/topic/%s/time-series/v2", topic)
	if err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Creators retrieves the top creators posting about a topic:
// GET /public/topic/:topic/creators/v1.
func (s *TopicsService) Creators(ctx context.Context, topic string) (*TopicCreatorsResponse, error) {
	var out TopicCreatorsResponse
	path := fmt.Sprintf("/public/topic/%s/creators/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// News retrieves news articles related to a topic:
// GET /public/topic/:topic/news/v1.
func (s *TopicsService) News(ctx context.Context, topic string) (*TopicNewsResponse, error) {
	var out TopicNewsResponse
	path := fmt.Sprintf("/public/topic/%s/news/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Posts retrieves social posts related to a topic:
// GET /public/topic/:topic/posts/v1.
func (s *TopicsService) Posts(ctx context.Context, topic string) (*TopicPostsResponse, error) {
	var out TopicPostsResponse
	path := fmt.Sprintf("/public/topic/%s/posts/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WhatsUp retrieves an AI-generated summary of recent activity for a
// topic: GET /public/topic/:topic/whatsup/v1.
func (s *TopicsService) WhatsUp(ctx context.Context, topic string) (*TopicWhatsUpResponse, error) {
	var out TopicWhatsUpResponse
	path := fmt.Sprintf("/public/topic/%s/whatsup/v1", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List retrieves the list of all trending topics:
// GET /public/topics/list/v1.
func (s *TopicsService) List(ctx context.Context) (*TopicsListResponse, error) {
	var out TopicsListResponse
	if err := s.client.doRequest(ctx, "GET", "/public/topics/list/v1", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
