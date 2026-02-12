package models

import (
	"time"
)

// Vip 模版
type Vip struct {
	Id int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null" json:"id"`
	Type string `gorm:"column:type;type:varchar(255);not null;comment:类型 " json:"type"` // 类型 
	IsEnable int8 `gorm:"column:is_enable;type:tinyint;not null;default:1;comment:是否启用" json:"is_enable"` // 是否启用
	Description string `gorm:"column:description;type:varchar(255);not null" json:"description"`
	OriginalPrice float64 `gorm:"column:original_price;type:decimal(10,2);not null;default:0.00;comment:单价 原价" json:"original_price"` // 单价 原价
	ActualPrice float64 `gorm:"column:actual_price;type:decimal(10,2);not null;default:0.00;comment:真实价格" json:"actual_price"` // 真实价格
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (Vip) TableName() string {
	return "vip"
}

