package dao

import (
	"context"
	"errors"
	"fmt"
	"golang-tars/internal/user/data/mysql/models"
	"golang-tars/pkg/db"

	"github.com/nicexiaonie/gconf"
	"gorm.io/gorm"
)

// UserDAO 用户数据访问接口
type UserDAO interface {
	// Create 创建用户
	Create(ctx context.Context, user *models.Users) error

	// FindByID 根据ID查询用户
	FindByID(ctx context.Context, id int64) (*models.Users, error)

	// FindByUsername 根据用户名查询用户
	FindByUsername(ctx context.Context, username string) (*models.Users, error)

	// FindByPhoneNumber 根据手机号查询用户
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Users, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *models.Users) error

	// Delete 删除用户（软删除）
	Delete(ctx context.Context, id int64) error

	// List 分页查询用户列表
	List(ctx context.Context, page, pageSize int) ([]*models.Users, int64, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByPhoneNumber 检查手机号是否存在
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
}

// userDAOImpl UserDAO 实现
type userDAOImpl struct {
	*db.BaseDAO
	instanceName string // 数据库实例名称
}

// NewUserDAO 创建 UserDAO 实例
// 这是 Wire Provider 函数，接受数据库管理器作为依赖
// 支持多数据库场景：可以从配置文件读取使用哪个数据库实例
func NewUserDAO(manager *db.Manager) (UserDAO, error) {
	// 从配置读取 User DAO 使用的数据库实例名称
	// 默认使用 "default"，也可以配置为其他实例
	instanceName := gconf.GetString("dao.user.db_instance")
	if instanceName == "" {
		instanceName = "default" // 默认使用 default 数据库
	}

	// 验证数据库实例是否存在
	if _, err := manager.GetClient(instanceName); err != nil {
		return nil, fmt.Errorf("failed to get database client '%s' for user dao: %w", instanceName, err)
	}

	// 使用 BaseDAO 提供的基础方法
	baseDAO, err := db.NewBaseDAO(instanceName, models.Users{}.TableName())
	if err != nil {
		return nil, fmt.Errorf("failed to create user dao: %w", err)
	}

	return &userDAOImpl{
		BaseDAO:      baseDAO,
		instanceName: instanceName,
	}, nil
}

// GetInstanceName 获取当前使用的数据库实例名称
func (d *userDAOImpl) GetInstanceName() string {
	return d.instanceName
}

// Create 创建用户
func (d *userDAOImpl) Create(ctx context.Context, user *models.Users) error {
	if err := d.BaseDAO.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByID 根据ID查询用户
func (d *userDAOImpl) FindByID(ctx context.Context, id int64) (*models.Users, error) {
	var user models.Users
	err := d.BaseDAO.FindByID(ctx, id, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: id=%d", id)
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (d *userDAOImpl) FindByUsername(ctx context.Context, username string) (*models.Users, error) {
	var user models.Users
	err := d.BaseDAO.FindOne(ctx, &user, "username = ?", username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: username=%s", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return &user, nil
}

// FindByPhoneNumber 根据手机号查询用户
func (d *userDAOImpl) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Users, error) {
	var user models.Users
	err := d.BaseDAO.FindOne(ctx, &user, "phone_number = ?", phoneNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: phone=%s", phoneNumber)
		}
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}
	return &user, nil
}

// Update 更新用户信息
func (d *userDAOImpl) Update(ctx context.Context, user *models.Users) error {
	if err := d.BaseDAO.Save(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 删除用户
func (d *userDAOImpl) Delete(ctx context.Context, id int64) error {
	if err := d.BaseDAO.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 分页查询用户列表
func (d *userDAOImpl) List(ctx context.Context, page, pageSize int) ([]*models.Users, int64, error) {
	var users []*models.Users

	result, err := d.BaseDAO.Paginate(ctx, &users, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, result.Total, nil
}

// ExistsByUsername 检查用户名是否存在
func (d *userDAOImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := d.BaseDAO.Exists(ctx, "username = ?", username)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return exists, nil
}

// ExistsByPhoneNumber 检查手机号是否存在
func (d *userDAOImpl) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	exists, err := d.BaseDAO.Exists(ctx, "phone_number = ?", phoneNumber)
	if err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}
	return exists, nil
}
