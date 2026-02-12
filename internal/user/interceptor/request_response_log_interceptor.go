package interceptor

import (
	"context"
	"encoding/json"
	"fmt"
	logger "golang-tars/pkg"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// RequestResponseLogInterceptor 请求响应参数日志拦截器
func RequestResponseLogInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure
			start := time.Now()

			// 获取日志对象
			log := logger.FromContext(ctx)

			// 序列化请求参数
			reqParam := serializeMessage(req.Any())

			log.Info(fmt.Sprintf("请求开始 procedure: %s, RequestBody: %s", procedure, reqParam))

			// 执行处理器
			resp, err := next(ctx, req)

			duration := time.Since(start)

			// 记录响应日志 - 美化格式
			if err != nil {
				// 错误情况：只记录错误信息（错误时响应体通常无效）
				log.Error(fmt.Sprintf("请求结束 procedure: %s, duration:%s, Error: %v", procedure, duration, err))
			} else {
				// 成功情况：安全地记录响应体
				respBody := serializeMessage(resp.Any())
				log.Info(fmt.Sprintf("请求结束 procedure: %s, duration:%s, ResponseBody: %s", procedure, duration, respBody))
			}

			return resp, err
		}
	}
}

// serializeMessage 序列化 protobuf 消息为 JSON 字符串
func serializeMessage(msg any) string {
	if msg == nil {
		return "{}"
	}

	// 尝试将消息转换为 proto.Message
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		// 如果不是 proto.Message，尝试直接 JSON 序列化
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Sprintf("error: %v", err)
		}
		return string(data)
	}

	// 使用 protojson 序列化 protobuf 消息
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true, // 支持0值字段
	}

	data, err := marshaler.Marshal(protoMsg)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}

	return string(data)
}
