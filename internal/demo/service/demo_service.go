package service

import (
	"context"
	"fmt"
	demov1 "golang-tars/pkg/proto/demo/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DemoService 实现 Demo gRPC 服务
type DemoService struct {
	demov1.UnimplementedDemoServiceServer
}

// NewDemoService 创建 Demo 服务实例
func NewDemoService() *DemoService {
	return &DemoService{}
}

// CreateDemo 创建 Demo
func (s *DemoService) CreateDemo(ctx context.Context, req *demov1.CreateDemoRequest) (*demov1.CreateDemoResponse, error) {
	// 参数验证（protoc-gen-validate 已自动生成验证代码）
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %v", err)
	}

	return &demov1.CreateDemoResponse{
		// Base: &commonv1.BaseResponse{
		// 	Code:    0,
		// 	Message: "success",
		// },
	}, nil
}

// GetDemo 获取 Demo
func (s *DemoService) GetDemo(ctx context.Context, req *demov1.GetDemoRequest) (*demov1.GetDemoResponse, error) {
	// 参数验证
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %v", err)
	}

	return &demov1.GetDemoResponse{}, nil
}

// UpdateDemo 更新 Demo
func (s *DemoService) UpdateDemo(ctx context.Context, req *demov1.UpdateDemoRequest) (*demov1.UpdateDemoResponse, error) {
	// 参数验证
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %v", err)
	}

	return &demov1.UpdateDemoResponse{}, nil
}

// DeleteDemo 删除 Demo
func (s *DemoService) DeleteDemo(ctx context.Context, req *demov1.DeleteDemoRequest) (*demov1.DeleteDemoResponse, error) {
	// 参数验证
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %v", err)
	}

	return &demov1.DeleteDemoResponse{
		// Success: true,
		Message: fmt.Sprintf("Demo %d 已成功删除", req.Id),
	}, nil
}

// ListDemo 列表 Demo（支持分页和筛选）
func (s *DemoService) ListDemo(ctx context.Context, req *demov1.ListDemoRequest) (*demov1.ListDemoResponse, error) {
	// 参数验证
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "参数验证失败: %v", err)
	}

	return &demov1.ListDemoResponse{}, nil
}

// contains 简单的字符串包含检查（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
		(len(s) > len(substr) && contains(s[1:], substr)))
}
