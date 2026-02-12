package models

import (
	"time"
)

// UserLoginLog 用户登录日志
type UserLoginLog struct {
	Id int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement;not null" json:"id"`
	UserId uint64 `gorm:"column:user_id;type:bigint unsigned;not null;comment:用户ID" json:"user_id"` // 用户ID
	LoginResule int8 `gorm:"column:login_resule;type:tinyint;not null;default:0;comment:登录结果" json:"login_resule"` // 登录结果
	LoginMode string `gorm:"column:login_mode;type:varchar(255);not null;comment:登录方式" json:"login_mode"` // 登录方式
	LoginIpV4 string `gorm:"column:login_ip_v4;type:varchar(32);not null;comment:IP" json:"login_ip_v4"` // IP
	LoginIpV6 string `gorm:"column:login_ip_v6;type:varchar(255);not null;comment:IP" json:"login_ip_v6"` // IP
	DeviceId string `gorm:"column:device_id;type:varchar(255);not null;comment:设备ID" json:"device_id"` // 设备ID
	DeviceFingerprint string `gorm:"column:device_fingerprint;type:varchar(255);not null;comment:设备指纹" json:"device_fingerprint"` // 设备指纹
	DeviceInfo string `gorm:"column:device_info;type:text;not null;comment:设备信息 json" json:"device_info"` // 设备信息 json
	Platform string `gorm:"column:platform;type:varchar(255);not null;comment:平台" json:"platform"` // 平台
	LoginTime time.Time `gorm:"column:login_time;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:登录时间" json:"login_time"` // 登录时间
	FailReason string `gorm:"column:fail_reason;type:varchar(255);not null;comment:失败原因" json:"fail_reason"` // 失败原因
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
}

// TableName 指定表名
func (UserLoginLog) TableName() string {
	return "user_login_log"
}

