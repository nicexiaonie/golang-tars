package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/nicexiaonie/glog"
	"github.com/sirupsen/logrus"
)

// ContextKey 自定义 context key 类型，避免与其他包的 key 冲突
type ContextKey string

const (
	// 追踪相关
	CtxKeyRequestID ContextKey = "RequestID"
	CtxKeyTraceID   ContextKey = "TraceID"
	CtxKeyLogger    ContextKey = "Logger"

	// 用户认证相关
	CtxKeyUserID            ContextKey = "UserID"
	CtxKeyUsername          ContextKey = "Username"
	CtxKeyPhoneNumber       ContextKey = "PhoneNumber"
	CtxKeyDeviceFingerprint ContextKey = "DeviceFingerprint"
	CtxKeyToken             ContextKey = "Token"
)

var Logger *logrus.Logger

// LoggerInit 初始化日志，返回错误供调用方处理
func LoggerInit(path string) error {
	var err error
	Logger, err = glog.New(&glog.Config{
		// 日志存放目录
		Path: path,
		// 日志文件名
		Filename: "app.log",
		// 日志级别 超过并包含设置级别以上的日志会处理保存
		Level: "info",
		// json text formatter
		Format: "json",
		// 自定义格式化接口
		//Formatter:    logrus.Formatter,

		// option: file
		Output:       "file",
		ReportCaller: true,

		// 可以为空，非空按照日期格式生成日志文件
		Split: ".2006010215",
		// 日志最后修改时间超过多少时间进行清理
		Lifetime: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	return nil
}

// WithFields 创建带字段的logger
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// FromContext 从context获取logger，如果不存在则返回全局Logger
func FromContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return Logger.WithFields(logrus.Fields{})
	}

	// 尝试从context获取logger
	if logger := GetLogger(ctx); logger != nil {
		if entry, ok := logger.(*logrus.Entry); ok {
			return entry
		}
	}

	// 如果context中没有logger，尝试获取request_id和trace_id
	requestID := GetRequestID(ctx)
	traceID := GetTraceID(ctx)

	fields := logrus.Fields{}
	if requestID != "" {
		fields["RequestID"] = requestID
	}
	if traceID != "" {
		fields["TraceID"] = traceID
	}

	return Logger.WithFields(fields)
}

// GetLogger 从context获取logger
func GetLogger(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}
	return ctx.Value(CtxKeyLogger)
}

// GetRequestID 从context获取请求ID
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// GetTraceID 从context获取追踪ID
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(CtxKeyTraceID).(string); ok {
		return v
	}
	return ""
}
