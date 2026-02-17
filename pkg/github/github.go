package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkglog "sb-scanner/pkg/logger"
)

type Client interface {
	SearchCommits(ctx context.Context, keywords []string, opts ...SearchOption) (SearchResult, error)
}

type DefaultClient struct {
	logger  *slog.Logger
	httpcli *http.Client
}

// Creates unauthenticated GitHub API client.
// Note: Unauthenticated requests are subject to much lower rate limits.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		logger:  pkglog.GetLogger().With("pkg", "github"),
		httpcli: http.DefaultClient,
	}
}

const (
	acceptHeaderValue     = "application/vnd.github+json"
	apiVersionHeaderKey   = "X-Github-Api-Version"
	apiVersionHeaderValue = "2022-11-28"
	userAgentHeaderValue  = "sb-scanner"
)

type SearchOption func(o *searchOptionValues)

type searchOptionValues struct {
	Page      *int
	Size      *int
	StartTime *time.Time
	EndTime   *time.Time

	authToken *string
}

func WithPage(page int) SearchOption {
	return func(o *searchOptionValues) {
		o.Page = &page
	}
}

func WithSize(size int) SearchOption {
	return func(o *searchOptionValues) {
		o.Size = &size
	}
}

func WithStartTime(t time.Time) SearchOption {
	return func(o *searchOptionValues) {
		o.StartTime = &t
	}
}

func WithEndTime(t time.Time) SearchOption {
	return func(o *searchOptionValues) {
		o.EndTime = &t
	}
}

// withAuthToken sets a personal access token for authentication; for internal use only
func withAuthToken(token string) SearchOption {
	return func(o *searchOptionValues) {
		o.authToken = &token
	}
}

func (c *DefaultClient) SearchCommits(ctx context.Context, keywords []string, opts ...SearchOption) (SearchResult, error) {
	if len(keywords) > 5 {
		return SearchResult{}, fmt.Errorf("maximum 5 keywords are allowed")
	}
	ov := &searchOptionValues{}
	for _, opt := range opts {
		opt(ov)
	}

	searchQ := strings.Join(keywords, " OR ")
	if ov.StartTime != nil {
		if ov.EndTime != nil {
			searchQ += fmt.Sprintf(" author-date:%s..%s", ov.StartTime.Format(time.RFC3339), ov.EndTime.Format(time.RFC3339))
		} else {
			searchQ += fmt.Sprintf(" author-date:%s..*", ov.StartTime.Format(time.RFC3339))
		}
	} else {
		if ov.EndTime != nil {
			searchQ += fmt.Sprintf(" author-date:*..%s", ov.EndTime.Format(time.RFC3339))
		}
	}
	c.logger.Debug("search query", "q", searchQ)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/search/commits", nil)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	q := req.URL.Query()
	q.Add("q", searchQ)
	q.Add("sort", "author-date")
	if ov.Page != nil {
		q.Add("page", strconv.Itoa(*ov.Page))
	}
	if ov.Size != nil {
		q.Add("per_page", strconv.Itoa(*ov.Size))
	}
	req.URL.RawQuery = q.Encode()
	c.logger.Debug("search request query string", "query", req.URL.RawQuery)
	req.Header.Add("Accept", acceptHeaderValue)
	req.Header.Add(apiVersionHeaderKey, apiVersionHeaderValue)
	req.Header.Add("User-Agent", userAgentHeaderValue)
	if ov.authToken != nil {
		req.Header.Add("Authorization", "Bearer "+*(ov.authToken))
	}

	resp, err := c.httpcli.Do(req)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to read response body: %w", err)
	}
	var result SearchResult
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return SearchResult{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return result, nil
}
