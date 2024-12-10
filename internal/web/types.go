package web

type Page struct {
	Page     int64 `json:"page" binding:"omitempty,gt=0"`
	PageSize int64 `json:"pageSize" binding:"omitempty,gt=0"`
}

func (p *Page) SetDefault() {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if p.Page == 0 {
		p.Page = 1
	}
}
