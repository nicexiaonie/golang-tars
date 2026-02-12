package dao

import (
	"context"
	"errors"
	"fmt"
	"go-base/internal/user/data/mysql/models"
	"go-base/pkg/db"

	"gorm.io/gorm"
)

// DemoDAO Demo 数据访问接口
type DemoDAO interface {
	// Create 创建 Demo
	Create(ctx context.Context, demo *models.Demo) error

	// FindByID 根据ID查询
	FindByID(ctx context.Context, id int64) (*models.Demo, error)

	// Update 更新 Demo
	Update(ctx context.Context, demo *models.Demo) error

	// Delete 删除 Demo
	Delete(ctx context.Context, id int64) error

	// List 分页查询列表
	List(ctx context.Context, page, pageSize int) ([]*models.Demo, int64, error)

	// FindByName 根据名称查询
	FindByName(ctx context.Context, name string) (*models.Demo, error)
}

// demoDAOImpl DemoDAO 实现
type demoDAOImpl struct {
	*db.BaseDAO
}

// NewDemoDAO 创建 DemoDAO 实例
// 这是 Wire Provider 函数，接受数据库管理器作为依赖
// 支持多数据库场景：可以从配置文件读取使用哪个数据库实例
func NewDemoDAO(manager *db.Manager) (DemoDAO, error) {
	// 创建 BaseDAO
	baseDAO, err := db.NewBaseDAO("default", models.Demo{}.TableName())
	if err != nil {
		return nil, fmt.Errorf("failed to create demo dao: %w", err)
	}

	return &demoDAOImpl{
		BaseDAO: baseDAO,
	}, nil
}

// Create 创建 Demo
func (d *demoDAOImpl) Create(ctx context.Context, demo *models.Demo) error {
	if err := d.BaseDAO.Create(ctx, demo); err != nil {
		return fmt.Errorf("failed to create demo: %w", err)
	}
	return nil
}

// FindByID 根据ID查询
func (d *demoDAOImpl) FindByID(ctx context.Context, id int64) (*models.Demo, error) {
	var demo models.Demo
	err := d.BaseDAO.FindByID(ctx, id, &demo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("demo not found: id=%d", id)
		}
		return nil, fmt.Errorf("failed to find demo: %w", err)
	}
	return &demo, nil
}

// Update 更新 Demo
func (d *demoDAOImpl) Update(ctx context.Context, demo *models.Demo) error {
	if err := d.BaseDAO.Save(ctx, demo); err != nil {
		return fmt.Errorf("failed to update demo: %w", err)
	}
	return nil
}

// Delete 删除 Demo
func (d *demoDAOImpl) Delete(ctx context.Context, id int64) error {
	if err := d.BaseDAO.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete demo: %w", err)
	}
	return nil
}

// List 分页查询列表
func (d *demoDAOImpl) List(ctx context.Context, page, pageSize int) ([]*models.Demo, int64, error) {
	var demos []*models.Demo
	result, err := d.BaseDAO.Paginate(ctx, &demos, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list demos: %w", err)
	}
	return demos, result.Total, nil
}

// FindByName 根据名称查询
func (d *demoDAOImpl) FindByName(ctx context.Context, name string) (*models.Demo, error) {
	var demo models.Demo
	err := d.BaseDAO.FindOne(ctx, &demo, "name = ?", name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("demo not found: name=%s", name)
		}
		return nil, fmt.Errorf("failed to find demo by name: %w", err)
	}
	return &demo, nil
}
