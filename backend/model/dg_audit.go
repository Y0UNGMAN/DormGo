package model

// AuditRequest 批量/单体审核请求参数
// 通常用于 Controller 层接收前端传来的 JSON 数据
type AuditRequest struct {
	Status string `json:"status"`
	IDs    []uint `json:"ids"`
}
