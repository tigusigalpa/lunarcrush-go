package lunarcrush

import (
	"context"
	"fmt"
)

// AIService provides access to the /public/ai* endpoints of the
// LunarCrush API, which return AI-generated summaries for topics and
// creators.
type AIService struct {
	client *Client
}

// AISummary represents an AI-generated natural-language summary as
// returned by the AI endpoints.
type AISummary struct {
	Subject string `json:"subject,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// AITopicResponse is the response envelope returned by
// GET /public/ai/topic/:topic.
type AITopicResponse struct {
	Data   AISummary              `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// AICreatorResponse is the response envelope returned by
// GET /public/ai/creator/:network/:id.
type AICreatorResponse struct {
	Data   AISummary              `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Topic retrieves an AI-generated summary for a topic:
// GET /public/ai/topic/:topic.
func (s *AIService) Topic(ctx context.Context, topic string) (*AITopicResponse, error) {
	var out AITopicResponse
	path := fmt.Sprintf("/public/ai/topic/%s", topic)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Creator retrieves an AI-generated summary for a creator:
// GET /public/ai/creator/:network/:id.
func (s *AIService) Creator(ctx context.Context, network, id string) (*AICreatorResponse, error) {
	var out AICreatorResponse
	path := fmt.Sprintf("/public/ai/creator/%s/%s", network, id)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
