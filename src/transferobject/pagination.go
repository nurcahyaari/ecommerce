package transferobject

import "github.com/nurcahyaari/ecommerce/src/domain/entity"

type Pagination struct {
	TotalPage int `json:"totalPage"`
	TotalSize int `json:"totalSize"`
	Page      int `json:"page"`
	Size      int `json:"size"`
}

func (p *Pagination) Default() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.Size == 0 {
		p.Size = 10
	}
}

func NewPagination(pagination entity.Pagination) Pagination {
	return Pagination{
		TotalPage: pagination.TotalPage,
		TotalSize: pagination.TotalSize,
		Page:      pagination.Page,
		Size:      pagination.Size,
	}
}
