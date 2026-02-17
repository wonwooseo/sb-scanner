package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"sb-scanner/model"
	"sb-scanner/pkg/github"
	pkglog "sb-scanner/pkg/logger"
	"sb-scanner/pkg/repository"
)

func Sync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Search and save commits.",
		Long:  "Search commits from GitHub and save to database.",
		Run: func(cmd *cobra.Command, args []string) {
			cfgF, err := cmd.Flags().GetString("config")
			if err != nil {
				slog.Error("failed to read config flag", "err", err)
				os.Exit(1)
			}
			v := viper.New()
			v.SetConfigFile(cfgF)
			if err := v.ReadInConfig(); err != nil {
				slog.Error("failed to read config", "err", err)
				os.Exit(1)
			}
			pkglog.InitLogger(v.GetString("loglevel"))
			logger := pkglog.GetLogger().With("cmd", "sync")

			var githubCli github.Client
			if v.GetBool("github.auth.enabled") {
				githubCli, err = github.NewAuthenticatedClient(v.GetString("github.auth.app_id"), v.GetString("github.auth.installation_id"), []byte(v.GetString("github.auth.key")))
				if err != nil {
					logger.Error("failed to initialize github client", "err", err)
					os.Exit(1)
				}
			} else {
				logger.Warn("github authentication is disabled, using unauthenticated client with lower rate limits")
				githubCli = github.NewDefaultClient()
			}
			repo, err := repository.NewRepository(v.GetString("db.url"), v.GetString("db.name"))
			if err != nil {
				logger.Error("failed to initialize repository", "err", err)
				os.Exit(1)
			}

			searchKeywords := v.GetStringSlice("github.search.keywords")
			if len(searchKeywords) == 0 {
				logger.Error("no search keywords provided in configuration")
				os.Exit(1)
			}
			perPage := v.GetInt("github.search.per_page")
			if perPage <= 0 || perPage > 100 {
				logger.Warn("invalid github.search.per_page, using default (30)")
				perPage = 30
			}
			maxPage := v.GetInt("github.search.max_pages")
			if maxPage <= 0 {
				logger.Warn("invalid github.search.max_pages, using default (10)")
				maxPage = 10
			}
			rateLimitWait := v.GetDuration("github.search.rate_limit_wait")
			if rateLimitWait < 0 {
				logger.Warn("invalid github.search.rate_limit_wait, using default (2s)")
				rateLimitWait = 2 * time.Second
			}

			now := time.Now()
			stimeF, _ := cmd.Flags().GetString("stime")
			stime, err := time.Parse(time.RFC3339, stimeF)
			if err != nil {
				logger.Warn("invalid stime format, using default value", "err", err)
				stime = now.Add(-24 * time.Hour)
			}
			etimeF, _ := cmd.Flags().GetString("etime")
			etime, err := time.Parse(time.RFC3339, etimeF)
			if err != nil {
				logger.Warn("invalid stime format, using default value", "err", err)
				etime = now
			}
			if etime.Before(stime) {
				logger.Warn("etime is before stime, adjusting etime to now")
				etime = now
			}

			h := &syncHandler{
				logger:         logger,
				githubCli:      githubCli,
				repo:           repo,
				searchPerPage:  perPage,
				searchMaxPage:  maxPage,
				rateLimitWait:  rateLimitWait,
				searchKeywords: searchKeywords,
			}
			if err := h.Run(stime, etime); err != nil {
				os.Exit(1)
			}
		},
	}

	flags := cmd.Flags()
	flags.String("stime", "", "search start time in RFC3339 foramt (default: 24 hours ago)")
	flags.String("etime", "", "search end time in RFC3339 format (default: now)")
	cmd.PersistentFlags().AddFlagSet(flags)
	return cmd
}

type syncHandler struct {
	logger    *slog.Logger
	githubCli github.Client
	repo      *repository.Repository

	searchPerPage  int
	searchMaxPage  int
	rateLimitWait  time.Duration
	searchKeywords []string
}

func (h *syncHandler) Run(stime, etime time.Time) error {
	// max 5 keywords per search due to GitHub Search API limitations
	for i := 0; i < len(h.searchKeywords); i += 5 {
		searchPage := 1
		for searchPage <= h.searchMaxPage {
			opts := []github.SearchOption{
				github.WithStartTime(stime),
				github.WithEndTime(etime),
				github.WithSize(h.searchPerPage),
				github.WithPage(searchPage),
			}
			searched, err := h.githubCli.SearchCommits(context.Background(), h.searchKeywords[i:min(i+5, len(h.searchKeywords))], opts...)
			if err != nil {
				h.logger.Error("failed to search commits", "err", err)
				return err
			}
			if len(searched.Items) == 0 {
				h.logger.Info("no more commits found, ending sync")
				break
			}
			h.logger.Debug("fetched commits from github", "page", searchPage, "items", len(searched.Items), "total_count", searched.TotalCount, "incomplete_results", searched.IncompleteResults)

			var commits []model.Commit
			for _, c := range searched.Items {
				commits = append(commits, model.Commit{
					ID:      fmt.Sprintf("%d:%s", c.Commit.Author.Date.Unix(), c.SHA[:7]),
					SHA:     c.SHA,
					URL:     c.HTMLURL,
					Message: c.Commit.Message,
					Author: model.Author{
						Username:  c.AuthorMeta.Login,
						AvatarURL: c.AuthorMeta.AvatarURL,
					},
					Time: c.Commit.Author.Date,
				})
			}
			inserted, err := h.repo.PutCommits(context.Background(), commits)
			if err != nil {
				h.logger.Error("failed to put commits to db", "err", err)
				return err
			}

			h.logger.Debug("inserted commits to database", "page", searchPage, "commits_found", len(commits), "commits_inserted", inserted)
			searchPage++

			if h.rateLimitWait > 0 {
				h.logger.Debug("waiting for rate limit", "duration", h.rateLimitWait.String())
				time.Sleep(h.rateLimitWait)
			}
		}
	}

	return nil
}
