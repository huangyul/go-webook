package domain

import "time"

type User struct {
	ID        int64
	Email     string
	Password  string
	AboutMe   string
	Birthday  time.Time
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
