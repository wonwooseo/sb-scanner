package repository

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sb-scanner/model"
	"sb-scanner/pkg/db/mongodb"
	pkglog "sb-scanner/pkg/logger"
)

type Repository struct {
	logger *slog.Logger
	dbcli  *mongo.Client

	database string
}

func NewRepository(dbURL, dbName string) (*Repository, error) {
	dbcli, err := mongodb.GetClient(dbURL)
	if err != nil {
		return nil, err
	}
	return &Repository{
		logger:   pkglog.GetLogger().With("pkg", "repository"),
		dbcli:    dbcli,
		database: dbName,
	}, nil
}

const collectionCommits = "commits"

func (r *Repository) PutCommits(ctx context.Context, commits []model.Commit) (int, error) {
	col := r.dbcli.Database(r.database).Collection(collectionCommits)

	input := make([]any, len(commits))
	for i, commit := range commits {
		doc := bson.M{
			"_id":     commit.ID,
			"sha":     commit.SHA,
			"url":     commit.URL,
			"message": commit.Message,
			"author":  commit.Author,
			"time":    commit.Time,
		}
		input[i] = doc
	}

	res, err := col.InsertMany(ctx, input, options.InsertMany().SetOrdered(false)) // ignore duplicates
	if err != nil {
		return 0, fmt.Errorf("failed to insert commit documents: %w", err)
	}

	return len(res.InsertedIDs), nil
}

func (r *Repository) GetCommits(ctx context.Context, bookmark *string, limit int64) ([]model.Commit, error) {
	col := r.dbcli.Database(r.database).Collection(collectionCommits)

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(limit) // time desc, limit
	filter := bson.M{}
	if bookmark != nil {
		filter["_id"] = bson.M{"$lt": *bookmark}
	}
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find commit documents from db: %w", err)
	}
	defer cursor.Close(ctx)

	var commits []model.Commit
	for cursor.Next(ctx) {
		var c model.Commit
		if err := cursor.Decode(&c); err != nil {
			return nil, fmt.Errorf("failed to decode document to go struct: %w", err)
		}
		commits = append(commits, c)
	}

	return commits, nil
}
