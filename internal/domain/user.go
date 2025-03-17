package domain

import "time"

type User struct {
	ID        int64
	Email     string
	Password  string
	Nickname  string
	Birthday  time.Time
	Phone     string
	AboutMe   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
