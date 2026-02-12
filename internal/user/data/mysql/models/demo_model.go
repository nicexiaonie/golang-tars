package models

import "time"

// Demo Demo 示例表
type Demo struct {
	ID          int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null;comment:名称" json:"name"`
	Description string    `gorm:"column:description;type:text;comment:描述" json:"description"`
	Status      int       `gorm:"column:status;type:int;not null;default:1;comment:状态(1:启用 0:禁用)" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (Demo) TableName() string {
	return "demo"
}
