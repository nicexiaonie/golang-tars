package response

// 错误码定义
const (
	// 通用错误码 (0-999)
	CodeSuccess      = 0    // 成功
	CodeError        = 1    // 通用错误
	CodeInvalidParam = 1001 // 参数错误
	CodeUnauthorized = 1002 // 未授权
	CodeForbidden    = 1003 // 禁止访问
	CodeNotFound     = 1004 // 资源不存在
	CodeTimeout      = 1005 // 请求超时
	CodeTooManyReq   = 1006 // 请求过于频繁

	// Token相关错误码 (1100-1199)
	CodeTokenMissing     = 1101 // Token缺失
	CodeTokenInvalid     = 1102 // Token无效
	CodeTokenExpired     = 1103 // Token已过期
	CodeTokenFormatError = 1104 // Token格式错误

	// 数据库错误 (2000-2999)
	CodeDBError       = 2001 // 数据库错误
	CodeDBNotFound    = 2002 // 数据不存在
	CodeDBDuplicate   = 2003 // 数据重复
	CodeDBTransaction = 2004 // 事务错误

	// 业务错误 - 用户相关 (3000-3999)
	CodeUserNotFound      = 3001 // 用户不存在
	CodeUserExists        = 3002 // 用户已存在
	CodeUserDisabled      = 3003 // 用户已禁用
	CodePasswordError     = 3004 // 密码错误
	CodePhoneExists       = 3005 // 手机号已存在
	CodePhoneNotFound     = 3006 // 手机号不存在
	CodeVerificationError = 3007 // 验证码错误
	CodeDeviceLimitReached = 3008 // 设备数量已达上限
	CodeDeviceOffline      = 3009 // 设备已离线，需要重新登录
	CodeHeartbeatFailed    = 3010 // 心跳上报失败

	// 业务错误 - 短信相关 (4000-4999)
	CodeSMSError      = 4001 // 短信发送失败
	CodeSMSTooFreq    = 4002 // 短信发送过于频繁
	CodeSMSVerifyFail = 4003 // 短信验证失败

	// 系统错误 (9000-9999)
	CodeSystemError   = 9001 // 系统错误
	CodeServiceError  = 9002 // 服务错误
	CodeNetworkError  = 9003 // 网络错误
	CodeUnknownError  = 9999 // 未知错误
)

// ErrorMessage 错误码对应的消息
var ErrorMessage = map[int]string{
	CodeSuccess:      "success",
	CodeError:        "error",
	CodeInvalidParam: "invalid parameter",
	CodeUnauthorized: "unauthorized",
	CodeForbidden:    "forbidden",
	CodeNotFound:     "not found",
	CodeTimeout:      "request timeout",
	CodeTooManyReq:   "too many requests",

	CodeTokenMissing:     "token missing",
	CodeTokenInvalid:     "token invalid",
	CodeTokenExpired:     "token expired",
	CodeTokenFormatError: "token format error",

	CodeDBError:       "database error",
	CodeDBNotFound:    "data not found",
	CodeDBDuplicate:   "data duplicate",
	CodeDBTransaction: "transaction error",

	CodeUserNotFound:      "user not found",
	CodeUserExists:        "user already exists",
	CodeUserDisabled:      "user disabled",
	CodePasswordError:     "password error",
	CodePhoneExists:       "phone number already exists",
	CodePhoneNotFound:     "phone number not found",
	CodeVerificationError: "verification code error",
	CodeDeviceLimitReached: "device limit reached",
	CodeDeviceOffline:      "device offline, please login again",
	CodeHeartbeatFailed:    "heartbeat failed",

	CodeSMSError:      "sms send failed",
	CodeSMSTooFreq:    "sms send too frequently",
	CodeSMSVerifyFail: "sms verification failed",

	CodeSystemError:  "system error",
	CodeServiceError: "service error",
	CodeNetworkError: "network error",
	CodeUnknownError: "unknown error",
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code int) string {
	if msg, ok := ErrorMessage[code]; ok {
		return msg
	}
	return ErrorMessage[CodeUnknownError]
}
