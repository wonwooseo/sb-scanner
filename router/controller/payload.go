package controller

import (
	"net/http"

	"sb-scanner/model"
)

type ResponseGetCommits struct {
	Commits  []model.Commit `json:"commits"`
	Bookmark *string        `json:"bookmark,omitempty"`
}

func (res *ResponseGetCommits) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
