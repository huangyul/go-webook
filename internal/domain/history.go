package domain

import "time"

type History struct {
	Id           int64
	ArticleId    int64
	ArticleTitle string
	AuthorId     int64
	AuthorName   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
