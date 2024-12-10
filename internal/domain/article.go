package domain

import "time"

type Article struct {
	ID        int64
	Title     string
	Content   string
	Status    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	Author
}

func (a Article) Abstract() string {
	str := []rune(a.Content)
	if len(str) > 128 {
		str = str[:128]
	}
	return string(str)
}

type Author struct {
	ID   int64
	Name string
}

type ArticleStatus uint8

func (s ArticleStatus) toUint8() uint8 {
	return uint8(s)
}

const (
	ArticleStatusUnknown = iota
	ArticleStatusUnPublished
	ArticleStatusPublished
	ArticleStatusWithdraw
)
