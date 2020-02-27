package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// For parsing API result

// Author _
type Author struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

// Commit _
type Commit struct {
	Author  *Author `json:"author"`
	Message string  `json:"message"`
}

// AuthorMeta _
type AuthorMeta struct {
	AvatarURL string `json:"avatar_url"`
}

// SearchResponseItem _
type SearchResponseItem struct {
	HTMLURL    string      `json:"html_url"`
	Commit     *Commit     `json:"commit"`
	AuthorMeta *AuthorMeta `json:"author"`
}

// SearchResponse _
type SearchResponse struct {
	Count int                   `json:"total_count"`
	Items []*SearchResponseItem `json:"items"`
}

// TrimmedCommitItem _
type TrimmedCommitItem struct {
	URL       string    `json:"url"`
	Author    string    `json:"author"`
	AvatarURL string    `json:"avatar_url"`
	Message   string    `json:"message"`
	Date      time.Time `json:"date"`
}

const apiURL = "https://api.github.com/search/commits"

var keywords = []string{"ㅅㅂ", "시바"} // Max. 5 OR operators

var client *http.Client

// SearchLatestCommit _
func SearchLatestCommit() ([]*TrimmedCommitItem, error) {
	if client == nil {
		log.Println("init new http client")
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}
	rawQuery := strings.Join(keywords, " OR ")
	encodedQuery := url.QueryEscape(rawQuery)
	searchURL := fmt.Sprintf("%s?q=%s&sort=author-date", apiURL, encodedQuery)
	log.Println("run new commit search")
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.cloak-preview")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, err
	}
	uResp := &SearchResponse{}
	if err := json.Unmarshal(bytes, uResp); err != nil {
		return nil, err
	}
	result := []*TrimmedCommitItem{}
	for _, item := range uResp.Items {
		trimmedItem := &TrimmedCommitItem{
			URL:     item.HTMLURL,
			Author:  item.Commit.Author.Name,
			Message: item.Commit.Message,
			Date:    item.Commit.Author.Date,
		}
		if item.AuthorMeta != nil {
			trimmedItem.AvatarURL = item.AuthorMeta.AvatarURL
		}
		result = append(result, trimmedItem)
	}
	log.Printf("finished new commit search (feched commits: %d)", len(result))
	return result, nil
}
