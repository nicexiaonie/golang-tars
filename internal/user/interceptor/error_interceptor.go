package interceptor

import (
	"context"
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
)

// UnifiedErrorResponse 统一错误响应格式
type UnifiedErrorResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorInterceptor 错误统一处理拦截器
// 将 Connect 错误转换为统一的 {code, message, data} 格式
func ErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// 执行处理器
			resp, err := next(ctx, req)

			// 如果没有错误，直接返回
			if err == nil {
				return resp, nil
			}

			// 如果有错误，将其转换为统一格式
			// 注意：这里我们仍然返回error，但通过HTTP状态码200和响应体来表达错误
			// 如果你想要HTTP状态码也是200，需要在HTTP层面处理
			return resp, err
		}
	}
}

// HTTPErrorHandler HTTP层的错误处理器
// 用于将所有错误响应转换为统一格式，并返回HTTP 200
func HTTPErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	// 从Connect错误中提取信息
	connectErr := connect.CodeOf(err)

	// 构建统一响应
	response := UnifiedErrorResponse{
		Code:    int32(connectErr),
		Message: err.Error(),
		Data:    nil,
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 始终返回200，错误信息在响应体中

	// 写入响应
	json.NewEncoder(w).Encode(response)
}
