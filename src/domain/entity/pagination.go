package entity

type Pagination struct {
	Page      int
	Size      int
	TotalPage int
	TotalSize int
}

func (p *Pagination) DefaultPagination() {
	p.Page = 1
	p.Size = 10
}

func NewPagination(totalRows, pageSize int) Pagination {
	totalPage := totalRows / pageSize
	return Pagination{
		TotalPage: totalPage,
		TotalSize: totalRows,
		Size:      pageSize,
	}
}

func (f Pagination) Pagination() (string, []interface{}) {
	page := (f.Page - 1) * f.Size
	args := make([]interface{}, 0)
	args = append(args, f.Size)
	args = append(args, page)
	return "LIMIT ? OFFSET ?", args
}
