package model

import "time"

type Commit struct {
	ID      string    `json:"-" bson:"_id"`
	SHA     string    `json:"sha" bson:"sha"`
	URL     string    `json:"url" bson:"url"`
	Message string    `json:"message" bson:"message"`
	Author  Author    `json:"author" bson:"author"`
	Time    time.Time `json:"time" bson:"time"`
}

type Author struct {
	Username  string `json:"username" bson:"username"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
}
