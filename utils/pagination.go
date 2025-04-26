package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// Pagination 结构体用于表示分页信息，包括页码、每页大小、偏移量和总数。
// 用于在数据库查询中实现分页功能。
type Pagination struct {
	Page   int // 页码，表示当前请求的页数
	Size   int // 每页大小，表示每页显示的记录数
	Offset int // 偏移量，计算分页时需要跳过的记录数
	Total  int // 总数，表示所有数据的总记录数
}

// NewPagination 创建一个新的分页对象。
// 该函数从 HTTP 请求的查询参数中获取分页信息（页码和每页大小），
// 如果没有提供分页信息，则使用默认值（页码为 1，每页大小为 10）。
// 参数:
//
//	c: Gin 上下文，用于获取查询参数。
//
// 返回:
//
//	返回一个指向 Pagination 结构体的指针。
func NewPagination(c *gin.Context) *Pagination {
	var p = &Pagination{}
	// 从查询参数中获取页码和每页大小
	p.Page, p.Size = cast.ToInt(c.Query("page")), cast.ToInt(c.Query("limit"))
	// 默认分页，如果页码或每页大小为 0，则使用默认值
	if p.Page == 0 || p.Size == 0 {
		p.Page, p.Size = 1, 10
	}
	// 计算偏移量
	p.Offset = (p.Page - 1) * p.Size
	return p
}

// Paginate 返回一个 GORM 查询的函数，
// 用于在数据库查询中应用分页限制，包括偏移量和每页大小。
//
// 返回:
//
//	返回一个接受 *gorm.DB 作为参数的函数，该函数会应用分页设置并返回一个分页后的查询。
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 使用 GORM 的 Offset 和 Limit 方法实现分页
		return db.Offset(p.Offset).Limit(p.Size)
	}
}
