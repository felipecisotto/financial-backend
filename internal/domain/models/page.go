package models

type Page[T any] struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"total_pages"`
	Results    []T   `json:"results"`
}

type PageRequest struct {
	Limit int64
	Page  int64
}

func (p *PageRequest) Offset() int {
	newPage := (int(p.Page) - 1)
	if (newPage < 0) {
		return 0
	}
	return newPage * int(p.Limit)
}
