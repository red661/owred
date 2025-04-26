package response

// Response 定义了通用响应的数据结构
type Response struct {
	Code    uint        `json:"code,omitempty"`    // 响应码
	Success bool        `json:"success"`           // 成功标志
	Message string      `json:"message,omitempty"` // 消息内容
	Data    interface{} `json:"data,omitempty"`    // 数据内容
}
