package model


const (
	Ok           int32 = 0
	ErrFormat    int32 = 1
	ErrNotFound  int32 = 2
	ErrForbidden int32 = 3
	ErrUnknown   int32 = 99
	ErrParam     int32 = 400
	ErrInternal  int32 = 500
)

// APIResult api调用结果
type APIResult struct {
	Code    int32    `json:"code"`
	Message string `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewApiResult 构造函数
func NewApiResult (code int32, msg string, data interface{}) APIResult {
	return APIResult{
		Code:code,
		Message: msg,
		Data: data,
	}
}

// NewOkResult 返回成功的 apiresult
func NewOkResult(data interface{}) APIResult {
	return APIResult{
		Code:    Ok,
		Message: "ok",
		Data:    data,
	}
}
