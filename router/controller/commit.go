package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	pkglog "sb-scanner/pkg/logger"
	"sb-scanner/pkg/repository"
)

type Controller struct {
	logger *slog.Logger
	repo   *repository.Repository
}

func NewController(repo *repository.Repository) *Controller {
	return &Controller{
		logger: pkglog.GetLogger().With("pkg", "controller"),
		repo:   repo,
	}
}

// GetCommits godoc
// @Summary      Gets list of commits
// @Tags         commit
// @Produce      json
// @Param        bookmark query string false "pagination bookmark"
// @Param        limit    query int    false "limit number of commits" default(30)
// @Success      200  {object}  ResponseGetCommits
// @Failure      500  {object}  rerr.ErrResponse
// @Router       /api/v1/commit [get]
func (c *Controller) GetCommits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := c.logger.With("func", "GetCommits")

	var bookmark *string
	if r.URL.Query().Get("bookmark") != "" {
		bookmark = new(string)
		*bookmark = r.URL.Query().Get("bookmark")
	}
	var limit int64 = 30
	if r.URL.Query().Get("limit") != "" {
		l, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			render.Render(w, r, MakeBadRequestError("query parameter `limit` must be an integer between 1 and 100"))
			return
		}
		if l > 100 || l <= 0 {
			render.Render(w, r, MakeBadRequestError("query parameter `limit` must be an integer between 1 and 100"))
			return
		}
		limit = l
	}

	commits, err := c.repo.GetCommits(ctx, bookmark, limit)
	if err != nil {
		logger.Error("failed to get commits from repository", "err", err)
		render.Render(w, r, MakeInternalServerError())
		return
	}
	var nextBookmark *string
	if int64(len(commits)) == limit {
		nextBookmarkCommit := commits[len(commits)-1]
		nextBookmark = &nextBookmarkCommit.ID
	}

	render.Render(w, r, &ResponseGetCommits{
		Commits:  commits,
		Bookmark: nextBookmark,
	})
}
