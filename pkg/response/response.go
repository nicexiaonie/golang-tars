package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(message string, data interface{}) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// WriteJSON 将响应写入 HTTP ResponseWriter
func (r *Response) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK) // HTTP 状态码始终返回 200，业务状态通过 code 字段区分
	if err := json.NewEncoder(w).Encode(r); err != nil {
		// WriteHeader 已发送，无法更改状态码，仅记录错误
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}

// ToJSON 转换为 JSON 字符串
func (r *Response) ToJSON() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
