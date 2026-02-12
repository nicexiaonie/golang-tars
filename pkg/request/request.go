package request

import (
	"encoding/json"
)

// Request 统一请求结构
// request_id和trace_id从HTTP Header中获取（X-Request-ID, X-Trace-ID）
type Request struct {
	Body json.RawMessage `json:"body"` // 业务数据（原始JSON）
}

// UnmarshalBody 解析业务数据到指定结构
func (r *Request) UnmarshalBody(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}
