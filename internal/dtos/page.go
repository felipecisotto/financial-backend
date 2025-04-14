package dtos

type PageRequest struct {
	Page  int64 `form:"page,default=1"`
	Limit int64 `form:"limit,default=10"`
}