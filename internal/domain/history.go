package domain

import "time"

type History struct {
	BizID     int64
	UserID    int64
	UpdatedAt time.Time
}
