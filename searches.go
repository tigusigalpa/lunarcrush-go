package lunarcrush

import (
	"context"
	"fmt"
)

// SearchesService provides access to the /public/searches* endpoints of
// the LunarCrush API, which allow creating and managing saved searches.
type SearchesService struct {
	client *Client
}

// Search represents a saved search as returned by the LunarCrush API.
type Search struct {
	Slug       string `json:"slug,omitempty"`
	Name       string `json:"name,omitempty"`
	SearchJSON string `json:"search_json,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}

// SearchResponse is the response envelope returned when creating,
// fetching, or updating a single search.
type SearchResponse struct {
	Data   Search                 `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// SearchListResponse is the response envelope returned by
// GET /public/searches/list.
type SearchListResponse struct {
	Data   []Search               `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// SearchDeleteResponse is the response envelope returned by
// GET /public/searches/:slug/delete.
type SearchDeleteResponse struct {
	Data   map[string]interface{} `json:"data"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// SearchCreateParams holds the parameters required to create a new
// saved search via GET /public/searches/create.
type SearchCreateParams struct {
	// Name is the unique name of the search. Required.
	Name string
	// SearchJSON is the JSON-encoded search definition. Required.
	SearchJSON string
	// Priority optionally sets the search's execution priority.
	Priority *int
}

// toQuery converts SearchCreateParams into a queryValues helper.
func (p *SearchCreateParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	name := p.Name
	searchJSON := p.SearchJSON
	q.setString("name", &name)
	q.setString("search_json", &searchJSON)
	q.setInt("priority", p.Priority)
	return q
}

// SearchUpdateParams holds the optional parameters accepted by the
// search update endpoint.
type SearchUpdateParams struct {
	// Name optionally renames the search.
	Name *string
	// SearchJSON optionally replaces the search definition.
	SearchJSON *string
	// Priority optionally updates the search's execution priority.
	Priority *int
}

// toQuery converts SearchUpdateParams into a queryValues helper.
func (p *SearchUpdateParams) toQuery() *queryValues {
	q := newQueryValues()
	if p == nil {
		return q
	}
	q.setString("name", p.Name)
	q.setString("search_json", p.SearchJSON)
	q.setInt("priority", p.Priority)
	return q
}

// Create creates a new saved search: GET /public/searches/create.
// Name and SearchJSON are required fields of params.
func (s *SearchesService) Create(ctx context.Context, params *SearchCreateParams) (*SearchResponse, error) {
	var out SearchResponse
	if params == nil {
		return nil, fmt.Errorf("lunarcrush: SearchCreateParams must not be nil")
	}
	if params.Name == "" || params.SearchJSON == "" {
		return nil, fmt.Errorf("lunarcrush: Name and SearchJSON are required")
	}
	err := s.client.doRequest(ctx, "GET", "/public/searches/create", params.toQuery().Values(), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// List retrieves all saved searches: GET /public/searches/list.
func (s *SearchesService) List(ctx context.Context) (*SearchListResponse, error) {
	var out SearchListResponse
	if err := s.client.doRequest(ctx, "GET", "/public/searches/list", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Search executes an ad-hoc search query: GET /public/searches/search.
func (s *SearchesService) Search(ctx context.Context, searchJSON string) (*SearchResponse, error) {
	var out SearchResponse
	q := newQueryValues()
	q.setString("search_json", &searchJSON)
	if err := s.client.doRequest(ctx, "GET", "/public/searches/search", q.Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a saved search by its slug: GET /public/searches/:slug.
func (s *SearchesService) Get(ctx context.Context, slug string) (*SearchResponse, error) {
	var out SearchResponse
	path := fmt.Sprintf("/public/searches/%s", slug)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update modifies an existing saved search:
// GET /public/searches/:slug/update.
func (s *SearchesService) Update(ctx context.Context, slug string, params *SearchUpdateParams) (*SearchResponse, error) {
	var out SearchResponse
	path := fmt.Sprintf("/public/searches/%s/update", slug)
	if err := s.client.doRequest(ctx, "GET", path, params.toQuery().Values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a saved search: GET /public/searches/:slug/delete.
func (s *SearchesService) Delete(ctx context.Context, slug string) (*SearchDeleteResponse, error) {
	var out SearchDeleteResponse
	path := fmt.Sprintf("/public/searches/%s/delete", slug)
	if err := s.client.doRequest(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
