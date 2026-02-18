package model

import "time"

type Commit struct {
	ID        string    `json:"-" bson:"_id"`
	SHA       string    `json:"sha" bson:"sha"`
	URL       string    `json:"url" bson:"url"`
	Message   string    `json:"message" bson:"message"`
	Author    Author    `json:"author" bson:"author"`
	Time      time.Time `json:"time" bson:"time"`
	Sentiment Sentiment `json:"sentiment" bson:"sentiment"`
}

type Author struct {
	Username  string `json:"username" bson:"username"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
}

type Sentiment struct {
	Score float64 `json:"score" bson:"score"` // -1.0 (negative) to 1.0 (positive)
	Model string  `json:"model" bson:"model"` // name of the model used for evaluation
}
