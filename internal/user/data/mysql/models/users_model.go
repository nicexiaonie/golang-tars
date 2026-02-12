package models

import (
	"time"
)

// Users 用户表
type Users struct {
	Id           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null" json:"id"`
	Username     string    `gorm:"column:username;type:varchar(255);not null;comment:账号" json:"username"`                                    // 账号
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(32);not null;comment:手机号码" json:"phone_number"`                           // 手机号码
	Nickname     string    `gorm:"column:nickname;type:varchar(255);not null;comment:用户昵称" json:"nickname"`                                  // 用户昵称
	Password     string    `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"password"`                                    // 密码
	WechatOpenId string    `gorm:"column:wechat_open_id;type:varchar(255);not null;comment:微信关联ID" json:"wechat_open_id"`                    // 微信关联ID
	RegisterTime time.Time `gorm:"column:register_time;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:注册时间" json:"reigster_time"` // 注册时间
	CreateTime   time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (Users) TableName() string {
	return "users"
}
