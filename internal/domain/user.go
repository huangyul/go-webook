package domain

import "time"

type User struct {
	ID int64

	Eamil    string
	Password string

	CreatedAt time.Time
	UpdatedAt time.Time
}
