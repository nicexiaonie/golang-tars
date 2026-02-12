package db

import (
	"context"
	"fmt"

	"github.com/nicexiaonie/grds"
	"gorm.io/gorm"
)

// BaseDAO 通用 DAO 基类
// 提供完整的增删改查、分页等基础方法
// 所有方法均接收 context.Context，支持超时取消传播
type BaseDAO struct {
	client       *grds.Client // grds 客户端
	tableName    string       // 表名
	instanceName string       // 数据库实例名称
}

// NewBaseDAO 创建 BaseDAO 实例
// instanceName: 数据库实例名称（如 "default"）
// tableName: 表名
func NewBaseDAO(instanceName, tableName string) (*BaseDAO, error) {
	// 从全局 Manager 获取客户端
	client, err := GetManager().GetClient(instanceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database client '%s': %w", instanceName, err)
	}

	return &BaseDAO{
		client:       client,
		tableName:    tableName,
		instanceName: instanceName,
	}, nil
}

// NewDefaultBaseDAO 创建使用默认数据库实例的 BaseDAO
// tableName: 表名
func NewDefaultBaseDAO(tableName string) (*BaseDAO, error) {
	return NewBaseDAO("default", tableName)
}

// GetClient 获取 grds 客户端
func (dao *BaseDAO) GetClient() *grds.Client {
	return dao.client
}

// DB 获取原始 gorm.DB 实例（不带 context，适用于不需要 context 的场景）
func (dao *BaseDAO) DB() *gorm.DB {
	return dao.client.DB()
}

// DBCtx 获取带有 context 的原始 gorm.DB 实例（不绑定表名，供复杂自定义查询使用）
func (dao *BaseDAO) DBCtx(ctx context.Context) *gorm.DB {
	return dao.client.DB().WithContext(ctx)
}

// db 内部方法：获取带 context + 表名的 gorm.DB（供 BaseDAO 内部 CRUD 方法使用）
func (dao *BaseDAO) db(ctx context.Context) *gorm.DB {
	return dao.client.DB().WithContext(ctx).Table(dao.tableName)
}

// Table 获取表查询构建器
func (dao *BaseDAO) Table() *grds.QueryBuilder {
	return dao.client.Table(dao.tableName)
}

// GetTableName 获取表名
func (dao *BaseDAO) GetTableName() string {
	return dao.tableName
}

// ==================== 创建方法 ====================

// Create 创建单条记录
func (dao *BaseDAO) Create(ctx context.Context, value interface{}) error {
	return dao.client.DB().WithContext(ctx).Create(value).Error
}

// CreateInBatches 批量创建记录
// batchSize: 每批次大小
func (dao *BaseDAO) CreateInBatches(ctx context.Context, values interface{}, batchSize int) error {
	return dao.db(ctx).CreateInBatches(values, batchSize).Error
}

// Save 保存记录（如果主键存在则更新，否则创建）
func (dao *BaseDAO) Save(ctx context.Context, value interface{}) error {
	return dao.client.DB().WithContext(ctx).Save(value).Error
}

// ==================== 更新方法 ====================

// Update 根据条件更新记录
// updates: 要更新的字段（map 或 struct）
// conditions: 更新条件（如 "id = ?", 1）
func (dao *BaseDAO) Update(ctx context.Context, updates interface{}, conditions ...interface{}) error {
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Updates(updates).Error
}

// UpdateByID 根据主键 ID 更新记录
func (dao *BaseDAO) UpdateByID(ctx context.Context, id interface{}, updates interface{}) error {
	return dao.db(ctx).Where("id = ?", id).Updates(updates).Error
}

// UpdateColumn 更新单个字段（包括零值）
func (dao *BaseDAO) UpdateColumn(ctx context.Context, column string, value interface{}, conditions ...interface{}) error {
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Update(column, value).Error
}

// ==================== 删除方法 ====================

// Delete 根据条件删除记录
func (dao *BaseDAO) Delete(ctx context.Context, conditions ...interface{}) error {
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Delete(nil).Error
}

// DeleteByID 根据主键 ID 删除记录
func (dao *BaseDAO) DeleteByID(ctx context.Context, id interface{}) error {
	return dao.db(ctx).Where("id = ?", id).Delete(nil).Error
}

// DeleteByIDs 根据主键 ID 列表批量删除
func (dao *BaseDAO) DeleteByIDs(ctx context.Context, ids []interface{}) error {
	return dao.db(ctx).Where("id IN ?", ids).Delete(nil).Error
}

// ==================== 查询方法 ====================

// FindByID 根据主键 ID 查询单条记录
func (dao *BaseDAO) FindByID(ctx context.Context, id interface{}, dest interface{}) error {
	return dao.db(ctx).Where("id = ?", id).First(dest).Error
}

// FindOne 根据条件查询单条记录
func (dao *BaseDAO) FindOne(ctx context.Context, dest interface{}, conditions ...interface{}) error {
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.First(dest).Error
}

// FindAll 查询所有记录
func (dao *BaseDAO) FindAll(ctx context.Context, dest interface{}, conditions ...interface{}) error {
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	return db.Find(dest).Error
}

// Count 统计记录数
func (dao *BaseDAO) Count(ctx context.Context, conditions ...interface{}) (int64, error) {
	var count int64
	db := dao.db(ctx)
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	err := db.Count(&count).Error
	return count, err
}

// Exists 判断记录是否存在
func (dao *BaseDAO) Exists(ctx context.Context, conditions ...interface{}) (bool, error) {
	count, err := dao.Count(ctx, conditions...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ==================== 分页方法 ====================

// Paginate 分页查询
func (dao *BaseDAO) Paginate(ctx context.Context, dest interface{}, page, pageSize int, scopes ...func(*gorm.DB) *gorm.DB) (*PageResult, error) {
	return Paginate(dao.db(ctx), dest, page, pageSize, scopes...)
}

// PaginateWithConditions 带条件的分页查询
func (dao *BaseDAO) PaginateWithConditions(ctx context.Context, dest interface{}, page, pageSize int, conditions ...interface{}) (*PageResult, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		if len(conditions) > 0 {
			return db.Where(conditions[0], conditions[1:]...)
		}
		return db
	}
	return Paginate(dao.db(ctx), dest, page, pageSize, scope)
}

// ==================== 事务方法 ====================

// Transaction 执行事务（支持 context 超时取消传播）
func (dao *BaseDAO) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return dao.client.DB().WithContext(ctx).Transaction(fn)
}

// ==================== 高级查询方法 ====================

// FindWithScope 使用自定义 scope 查询
func (dao *BaseDAO) FindWithScope(ctx context.Context, dest interface{}, scopes ...func(*gorm.DB) *gorm.DB) error {
	db := dao.db(ctx)
	for _, scope := range scopes {
		db = scope(db)
	}
	return db.Find(dest).Error
}

// FirstWithScope 使用自定义 scope 查询单条记录
func (dao *BaseDAO) FirstWithScope(ctx context.Context, dest interface{}, scopes ...func(*gorm.DB) *gorm.DB) error {
	db := dao.db(ctx)
	for _, scope := range scopes {
		db = scope(db)
	}
	return db.First(dest).Error
}
