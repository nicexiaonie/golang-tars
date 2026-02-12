package dao

import (
	"context"
	"fmt"
	"go-base/internal/user/data/mysql/models"
	"go-base/pkg/db"
)

// UserLoginLogDAO 用户登录日志 DAO
type UserLoginLogDAO struct {
	*db.BaseDAO
}

// NewUserLoginLogDAO 创建用户登录日志 DAO 实例
func NewUserLoginLogDAO() (*UserLoginLogDAO, error) {
	baseDAO, err := db.NewDefaultBaseDAO(models.UserLoginLog{}.TableName())
	if err != nil {
		return nil, err
	}
	return &UserLoginLogDAO{BaseDAO: baseDAO}, nil
}

// CreateLoginLog 创建登录日志
func (dao *UserLoginLogDAO) CreateLoginLog(ctx context.Context, log *models.UserLoginLog) error {
	if err := dao.BaseDAO.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create login log: %w", err)
	}
	return nil
}

// FindByUserID 根据用户ID查询登录日志
func (dao *UserLoginLogDAO) FindByUserID(ctx context.Context, userID int64, page, pageSize int) (*db.PageResult, error) {
	var logs []models.UserLoginLog
	result, err := dao.BaseDAO.PaginateWithConditions(ctx, &logs, page, pageSize, "user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find login logs: %w", err)
	}
	return result, nil
}
