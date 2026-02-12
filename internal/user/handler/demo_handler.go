package handler

import (
	"context"
	"fmt"
	"go-base/internal/user/service"
	demov1 "go-base/pkg/proto/demo/v1"

	"connectrpc.com/connect"
)

// DemoHandler Connect Demo 处理器（适配层：将 gRPC/Connect 请求转换为 Service 调用）
type DemoHandler struct {
	demoService service.DemoService
}

// NewDemoHandler 创建 Demo Handler 实例
func NewDemoHandler(demoService service.DemoService) *DemoHandler {
	return &DemoHandler{
		demoService: demoService,
	}
}

// CreateDemo 创建 Demo（Connect Handler）
func (h *DemoHandler) CreateDemo(
	ctx context.Context,
	req *connect.Request[demov1.CreateDemoRequest],
) (*connect.Response[demov1.CreateDemoResponse], error) {
	// 参数验证（使用 protobuf 生成的验证）
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数验证失败: %w", err))
	}

	// 调用 Service 层
	resp, err := h.demoService.CreateDemo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("创建失败: %w", err))
	}

	return connect.NewResponse(resp), nil
}

// GetDemo 获取 Demo（Connect Handler）
func (h *DemoHandler) GetDemo(
	ctx context.Context,
	req *connect.Request[demov1.GetDemoRequest],
) (*connect.Response[demov1.GetDemoResponse], error) {
	// 参数验证
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数验证失败: %w", err))
	}

	// 调用 Service 层
	resp, err := h.demoService.GetDemo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("获取失败: %w", err))
	}

	return connect.NewResponse(resp), nil
}

// UpdateDemo 更新 Demo（Connect Handler）
func (h *DemoHandler) UpdateDemo(
	ctx context.Context,
	req *connect.Request[demov1.UpdateDemoRequest],
) (*connect.Response[demov1.UpdateDemoResponse], error) {
	// 参数验证
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数验证失败: %w", err))
	}

	// 调用 Service 层
	resp, err := h.demoService.UpdateDemo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("更新失败: %w", err))
	}
	return connect.NewResponse(resp), nil
}

// DeleteDemo 删除 Demo（Connect Handler）
func (h *DemoHandler) DeleteDemo(
	ctx context.Context,
	req *connect.Request[demov1.DeleteDemoRequest],
) (*connect.Response[demov1.DeleteDemoResponse], error) {
	// 参数验证
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数验证失败: %w", err))
	}

	// 调用 Service 层
	resp, err := h.demoService.DeleteDemo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("删除失败: %w", err))
	}

	return connect.NewResponse(resp), nil
}

// ListDemo 列表查询 Demo（Connect Handler）
func (h *DemoHandler) ListDemo(
	ctx context.Context,
	req *connect.Request[demov1.ListDemoRequest],
) (*connect.Response[demov1.ListDemoResponse], error) {
	// 参数验证
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("参数验证失败: %w", err))
	}

	// 调用 Service 层
	resp, err := h.demoService.ListDemo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("查询失败: %w", err))
	}

	return connect.NewResponse(resp), nil
}
