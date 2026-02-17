package github

import "time"

type accessTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Author struct {
	Date time.Time `json:"date"`
}

type Commit struct {
	Author  Author `json:"author"`
	Message string `json:"message"`
}

type AuthorMeta struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type SearchResultItem struct {
	SHA        string     `json:"sha"`
	HTMLURL    string     `json:"html_url"`
	Commit     Commit     `json:"commit"`
	AuthorMeta AuthorMeta `json:"author"`
}

type SearchResult struct {
	TotalCount        int                `json:"total_count"`
	IncompleteResults bool               `json:"incomplete_results"`
	Items             []SearchResultItem `json:"items"`
}
