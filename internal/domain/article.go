package domain

import "time"

type Article struct {
	ID      int64
	Title   string
	Content string
	Author
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Author struct {
	ID   int64
	Name string
}
