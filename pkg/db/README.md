# 数据库管理模块

## 设计理念

**最小化封装，直接使用 grds**

- ✅ 不重复封装 grds 已有的功能
- ✅ Manager 只负责管理多个 grds.Client 实例
- ✅ 提供分页等常用辅助函数

## 快速开始

### 1. 初始化数据库

```go
import "cue-words/pkg/db"

dbManager := db.GetManager()

dbConfig := &db.DBConfig{
    Host:         "localhost",
    Port:         3306,
    Username:     "root",
    Password:     "password",
    Database:     "mydb",
    MaxIdleConns: 100,
    MaxOpenConns: 200,
}

err := dbManager.InitDB("default", dbConfig)
```

### 2. 获取 grds.Client 并使用

```go
// 获取 grds.Client
client, _ := dbManager.GetClient("default")

// 直接使用 grds 的方法
user := &models.Users{Username: "test"}
client.Create(user)

// 查询
var users []models.Users
client.Find(&users)

// 条件查询
client.Where("status = ?", 1).Find(&users)
```

### 3. 使用事务（grds 提供）

```go
client, _ := dbManager.GetClient("default")

// 使用 grds 的事务方法
err := client.Transaction(func(tx *gorm.DB) error {
    user1 := &models.Users{Username: "user1"}
    if err := tx.Create(user1).Error; err != nil {
        return err
    }

    user2 := &models.Users{Username: "user2"}
    if err := tx.Create(user2).Error; err != nil {
        return err
    }

    return nil // 返回 nil 自动提交
})
```

### 4. 使用分页（pkg/db 提供的辅助函数）

```go
import "cue-words/pkg/db"

client, _ := dbManager.GetClient("default")

var users []models.Users
result, err := db.Paginate(client.DB(), &users, 1, 10, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", 1)
})

// 使用分页结果
fmt.Println(result.Total)      // 总记录数
fmt.Println(result.TotalPages) // 总页数
```

## 多数据库管理

```go
dbManager := db.GetManager()

// 初始化多个数据库
dbManager.InitDB("main", &db.DBConfig{...})
dbManager.InitDB("analytics", &db.DBConfig{...})

// 使用不同的数据库
mainClient, _ := dbManager.GetClient("main")
analyticsClient, _ := dbManager.GetClient("analytics")

mainClient.Create(&user)
analyticsClient.Create(&log)
```

## DAO 层设计

**直接使用 grds.Client，无需工厂：**

```go
// 1. 获取 grds.Client
client, _ := db.GetManager().GetClient("default")

// 2. 创建 DAO
userDAO := dao.NewUserDAO(client)

// 3. 使用 DAO 的业务方法
user, _ := userDAO.FindByUsername("test")

// 4. 或直接使用 grds.Client
client.Create(&user)
client.Find(&users)
```

**DAO 实现示例：**

```go
type UserDAO struct {
    client *grds.Client
    table  string
}

func NewUserDAO(client *grds.Client) *UserDAO {
    return &UserDAO{
        client: client,
        table:  "users",
    }
}

// 业务方法
func (dao *UserDAO) FindByUsername(username string) (*models.Users, error) {
    var user models.Users
    err := dao.client.Table(dao.table).
        Where("username = ?", username).
        First(&user)
    return &user, err
}
```

## Manager API

Manager 只提供最小化的管理方法：

```go
// 初始化数据库
InitDB(name string, config *DBConfig) error

// 获取 grds.Client
GetClient(name string) (*grds.Client, error)
GetDefaultClient() (*grds.Client, error)

// 关闭连接
Close(name string) error
CloseAll() error

// 列出所有实例
ListInstances() []string
```

## grds.Client 常用方法

直接使用 grds 提供的方法：

```go
// CRUD
client.Create(&user)
client.Find(&users)
client.First(&user, "id = ?", 1)
client.Save(&user)
client.Delete(&user)

// 查询构建器
client.Table("users").Where("status = ?", 1).Find(&users)

// 事务
client.Transaction(func(tx *gorm.DB) error { ... })

// 健康检查
client.HealthCheck()
client.StatsInfo()
```

## 总结

**核心原则：**
1. ✅ Manager 只管理多数据库实例
2. ✅ 直接使用 grds.Client 的方法
3. ✅ 不重复封装已有功能
4. ✅ 保持简单和灵活

**参考文档：**
- [grds 项目地址](https://github.com/nicexiaonie/grds)
- [GORM 官方文档](https://gorm.io/docs/)
