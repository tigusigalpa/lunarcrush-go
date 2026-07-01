package lunarcrush

import (
	"context"
)

// SystemService provides access to the /public/system* endpoints of the
// LunarCrush API.
type SystemService struct {
	client *Client
}

// SystemChange represents a single changelog entry describing a change
// to the LunarCrush API or platform.
type SystemChange struct {
	Category    string `json:"category,omitempty"`
	Change      string `json:"change,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   int64  `json:"timestamp,omitempty"`
}

// SystemChangesResponse is the response envelope returned by
// GET /public/system/changes.
type SystemChangesResponse struct {
	Data   []SystemChange         `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// Changes retrieves the list of recent system/API changes:
// GET /public/system/changes.
func (s *SystemService) Changes(ctx context.Context) (*SystemChangesResponse, error) {
	var out SystemChangesResponse
	if err := s.client.doRequest(ctx, "GET", "/public/system/changes", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
