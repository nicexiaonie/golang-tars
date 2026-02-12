package db

import (
	"gorm.io/gorm"
)

// PageResult 分页结果
type PageResult struct {
	List       interface{} `json:"list"`        // 数据列表
	Total      int64       `json:"total"`       // 总记录数
	Page       int         `json:"page"`        // 当前页码
	PageSize   int         `json:"page_size"`   // 每页大小
	TotalPages int         `json:"total_pages"` // 总页数
}

// Paginate 分页查询辅助函数
// 使用示例：
//
//	var users []models.Users
//	result, err := db.Paginate(client.DB(), &users, 1, 10, func(db *gorm.DB) *gorm.DB {
//	    return db.Where("status = ?", 1)
//	})
func Paginate(db *gorm.DB, dest interface{}, page, pageSize int, scopes ...func(*gorm.DB) *gorm.DB) (*PageResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var total int64

	// 应用查询条件
	query := db
	for _, scope := range scopes {
		query = scope(query)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult{
		List:       dest,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
