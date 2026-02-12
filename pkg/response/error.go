package response

import "fmt"

// BizError 业务错误
type BizError struct {
	Code    int         // 错误码
	Message string      // 错误消息
	Data    interface{} // 附加数据
}

// Error 实现 error 接口
func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewBizError 创建业务错误
func NewBizError(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// NewBizErrorWithData 创建业务错误（带数据）
func NewBizErrorWithData(code int, message string, data interface{}) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewBizErrorByCode 根据错误码创建业务错误
func NewBizErrorByCode(code int) *BizError {
	return &BizError{
		Code:    code,
		Message: GetErrorMessage(code),
	}
}

// IsBizError 判断是否为业务错误
func IsBizError(err error) bool {
	_, ok := err.(*BizError)
	return ok
}

// GetBizError 获取业务错误
func GetBizError(err error) *BizError {
	if bizErr, ok := err.(*BizError); ok {
		return bizErr
	}
	return nil
}
