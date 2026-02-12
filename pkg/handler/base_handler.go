package handler

import (
	"encoding/json"
	"fmt"
	logger "go-base/pkg"
	"go-base/pkg/response"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// BaseHandler 基础处理器，提供通用方法
type BaseHandler struct{}

// HandleError 统一的错误处理
func (h *BaseHandler) HandleError(w http.ResponseWriter, err error) {
	// 判断是否为业务错误
	if bizErr := response.GetBizError(err); bizErr != nil {
		// 业务错误，返回错误码和消息
		response.ErrorWithData(bizErr.Code, bizErr.Message, bizErr.Data).WriteJSON(w)
		return
	}

	// 系统错误，记录日志并返回通用错误
	logger.Logger.Error(fmt.Sprintf("system error: %v", err))
	response.Error(response.CodeSystemError, "system error").WriteJSON(w)
}

// Success 返回成功响应
func (h *BaseHandler) Success(w http.ResponseWriter, data interface{}) {
	response.Success(data).WriteJSON(w)
}

// SuccessWithMessage 返回成功响应（自定义消息）
func (h *BaseHandler) SuccessWithMessage(w http.ResponseWriter, message string, data interface{}) {
	response.SuccessWithMessage(message, data).WriteJSON(w)
}

// Error 返回错误响应
func (h *BaseHandler) Error(w http.ResponseWriter, code int, message string) {
	response.Error(code, message).WriteJSON(w)
}

// ParseJSON 解析JSON请求体
func (h *BaseHandler) ParseJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to read request body: %v", err))
		return response.NewBizError(response.CodeInvalidParam, "failed to read request body")
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, v)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to unmarshal request: %v", err))
		return response.NewBizError(response.CodeInvalidParam, "invalid request format")
	}

	return nil
}

// GetQueryString 获取字符串查询参数
func (h *BaseHandler) GetQueryString(r *http.Request, key string, defaultValue ...string) string {
	value := r.URL.Query().Get(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// GetQueryInt 获取整数查询参数
func (h *BaseHandler) GetQueryInt(r *http.Request, key string, defaultValue ...int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, response.NewBizError(response.CodeInvalidParam, fmt.Sprintf("invalid %s format", key))
	}
	return intValue, nil
}

// GetQueryInt64 获取int64查询参数
func (h *BaseHandler) GetQueryInt64(r *http.Request, key string, defaultValue ...int64) (int64, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, nil
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, response.NewBizError(response.CodeInvalidParam, fmt.Sprintf("invalid %s format", key))
	}
	return int64Value, nil
}

// GetClientIP 获取客户端IP地址
func (h *BaseHandler) GetClientIP(r *http.Request) string {
	// 优先从X-Forwarded-For获取（如果使用了代理）
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		if idx := strings.IndexByte(ip, ','); idx > 0 {
			ip = ip[:idx]
		}
		return ip
	}

	// 其次从X-Real-IP获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 最后从RemoteAddr获取
	ip = r.RemoteAddr
	// RemoteAddr格式为 "IP:Port"，需要去掉端口
	if idx := strings.IndexByte(ip, ':'); idx > 0 {
		ip = ip[:idx]
	}
	return ip
}
