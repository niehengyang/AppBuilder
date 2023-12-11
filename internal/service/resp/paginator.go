package resp

import (
	"gorm.io/gorm"
	"math"
)

// PageOption 分页参数
type PageOption struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

type pageMeta struct {
	TotalRecord int64 `json:"total"`
	TotalPage   int   `json:"total_pages"`
	Offset      int   `json:"offset"`
	Limit       int   `json:"per_page"`
	Page        int   `json:"current_page"`
}

// Paginator 分页返回
type Paginator struct {
	Pagination pageMeta    `json:"pagination"`
	Records    interface{} `json:"data"`
}

// Paging 分页
func paging(p *PageOption, result interface{}) *Paginator {

	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	var paginator Paginator
	var count int64
	var offset int

	db.Model(result).Count(&count)
	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	//<-done

	paginator.Records = result
	paginator.Pagination.TotalRecord = count
	paginator.Pagination.Page = p.Page
	paginator.Pagination.Offset = offset
	paginator.Pagination.Limit = p.Limit
	paginator.Pagination.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))
	return &paginator
}

// Pagination
//
//	@Description: 分页
//	@param option 分页参数
//	@param result  分页数据
//	@param args 数据过滤所需的额外参数
//	@param callback 数据过滤
//	@return *types.Response
func Pagination(option *PageOption, result interface{}, callback func(result interface{}, arg []interface{}) interface{}, args ...interface{}) *Response {
	paginator := paging(option, result)
	paginator.Records = callback(paginator.Records, args)
	return &Response{
		Message: "success",
		Data:    paginator,
	}
}

func AfterPaginator(result interface{}, arg []interface{}) interface{} {
	return result
}
