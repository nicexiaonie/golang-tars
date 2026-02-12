package service

import (
	"context"
	"fmt"
	"go-base/internal/user/data/mysql/dao"
	"go-base/internal/user/data/mysql/models"
	"go-base/internal/user/proxy"
	logger "go-base/pkg"
	demov1 "go-base/pkg/proto/demo/v1"
)

// DemoService Demo 业务逻辑接口
type DemoService interface {
	CreateDemo(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.CreateDemoResponse, error)
	GetDemo(ctx context.Context, req *demov1.GetDemoRequest) (*demov1.GetDemoResponse, error)
	UpdateDemo(ctx context.Context, req *demov1.UpdateDemoRequest) (*demov1.UpdateDemoResponse, error)
	DeleteDemo(ctx context.Context, req *demov1.DeleteDemoRequest) (*demov1.DeleteDemoResponse, error)
	ListDemo(ctx context.Context, req *demov1.ListDemoRequest) (*demov1.ListDemoResponse, error)
}

// demoServiceImpl Demo Service 实现（集成 DAO 和 Proxy）
type demoServiceImpl struct {
	demoDAO   dao.DemoDAO     // 数据库访问层
	demoProxy proxy.DemoProxy // 第三方服务代理
}

// NewDemoService 创建 Demo Service 实例
// Wire 会自动注入 demoDAO 和 demoProxy 依赖
func NewDemoService(demoDAO dao.DemoDAO, demoProxy proxy.DemoProxy) DemoService {
	return &demoServiceImpl{
		demoDAO:   demoDAO,
		demoProxy: demoProxy,
	}
}

// CreateDemo 创建 Demo
func (s *demoServiceImpl) CreateDemo(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.CreateDemoResponse, error) {
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("Creating demo: name=%s, status=%d", req.Name, req.Status))

	// 1. 调用第三方服务验证（可选）
	// externalResp, err := s.demoProxy.GetUserInfo(ctx, req)
	// if err != nil {
	// 	log.Error(fmt.Sprintf("External service call failed: %v", err))
	// 	// 决定是否继续或返回错误
	// }

	// 2. 保存到数据库
	demo := &models.Demo{
		Name:        req.Name,
		Description: req.Description,
		Status:      int(req.Status),
	}

	if err := s.demoDAO.Create(ctx, demo); err != nil {
		log.Error(fmt.Sprintf("Failed to create demo: %v", err))
		return &demov1.CreateDemoResponse{
			Code:    500,
			Message: "创建失败",
		}, err
	}

	log.Info(fmt.Sprintf("Demo created successfully: id=%d", demo.ID))

	return &demov1.CreateDemoResponse{
		Code:    0,
		Message: "success",
		Data: &demov1.Demo{
			Id:          demo.ID,
			Name:        demo.Name,
			Description: demo.Description,
			Status:      int32(demo.Status),
		},
	}, nil
}

// GetDemo 获取 Demo
func (s *demoServiceImpl) GetDemo(ctx context.Context, req *demov1.GetDemoRequest) (*demov1.GetDemoResponse, error) {
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("Getting demo: id=%d", req.Id))

	// 从数据库查询
	demo, err := s.demoDAO.FindByID(ctx, req.Id)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get demo: %v", err))
		return &demov1.GetDemoResponse{
			Code:    404,
			Message: "Demo 不存在",
		}, err
	}

	return &demov1.GetDemoResponse{
		Code:    0,
		Message: "success",
		Data: &demov1.Demo{
			Id:          demo.ID,
			Name:        demo.Name,
			Description: demo.Description,
			Status:      int32(demo.Status),
		},
	}, nil
}

// UpdateDemo 更新 Demo
func (s *demoServiceImpl) UpdateDemo(ctx context.Context, req *demov1.UpdateDemoRequest) (*demov1.UpdateDemoResponse, error) {
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("Updating demo: id=%d", req.Id))

	// 1. 先查询是否存在
	demo, err := s.demoDAO.FindByID(ctx, req.Id)
	if err != nil {
		log.Error(fmt.Sprintf("Demo not found: %v", err))
		return &demov1.UpdateDemoResponse{
			Code:    404,
			Message: "Demo 不存在",
		}, err
	}

	// 2. 更新字段（处理可选字段）
	if req.Name != nil {
		demo.Name = *req.Name
	}
	if req.Description != nil {
		demo.Description = *req.Description
	}
	if req.Status != nil {
		demo.Status = int(*req.Status)
	}

	// 3. 保存到数据库
	if err := s.demoDAO.Update(ctx, demo); err != nil {
		log.Error(fmt.Sprintf("Failed to update demo: %v", err))
		return &demov1.UpdateDemoResponse{
			Code:    500,
			Message: "更新失败",
		}, err
	}

	log.Info(fmt.Sprintf("Demo updated successfully: id=%d", demo.ID))

	return &demov1.UpdateDemoResponse{
		Code:    0,
		Message: "success",
		Data: &demov1.Demo{
			Id:          demo.ID,
			Name:        demo.Name,
			Description: demo.Description,
			Status:      int32(demo.Status),
		},
	}, nil
}

// DeleteDemo 删除 Demo
func (s *demoServiceImpl) DeleteDemo(ctx context.Context, req *demov1.DeleteDemoRequest) (*demov1.DeleteDemoResponse, error) {
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("Deleting demo: id=%d", req.Id))

	// 删除记录
	if err := s.demoDAO.Delete(ctx, req.Id); err != nil {
		log.Error(fmt.Sprintf("Failed to delete demo: %v", err))
		return &demov1.DeleteDemoResponse{
			Code:    500,
			Message: "删除失败",
		}, err
	}

	log.Info(fmt.Sprintf("Demo deleted successfully: id=%d", req.Id))

	return &demov1.DeleteDemoResponse{
		Code:    0,
		Message: "success",
	}, nil
}

// ListDemo 列表查询 Demo
func (s *demoServiceImpl) ListDemo(ctx context.Context, req *demov1.ListDemoRequest) (*demov1.ListDemoResponse, error) {
	log := logger.FromContext(ctx)
	log.Info(fmt.Sprintf("Listing demos: page=%d, pageSize=%d", req.Page, req.PageSize))

	// 分页查询
	demos, total, err := s.demoDAO.List(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		log.Error(fmt.Sprintf("Failed to list demos: %v", err))
		return &demov1.ListDemoResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	// 转换为响应格式
	var demoList []*demov1.Demo
	for _, demo := range demos {
		demoList = append(demoList, &demov1.Demo{
			Id:          demo.ID,
			Name:        demo.Name,
			Description: demo.Description,
			Status:      int32(demo.Status),
		})
	}

	// 计算总页数
	totalPages := int32(total / int64(req.PageSize))
	if total%int64(req.PageSize) > 0 {
		totalPages++
	}

	return &demov1.ListDemoResponse{
		Code:    0,
		Message: "success",
		Data: &demov1.ListDemoData{
			Items:      demoList,
			Page:       req.Page,
			PageSize:   req.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
