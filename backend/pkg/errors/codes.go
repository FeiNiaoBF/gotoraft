package errors

// 错误码定义
const (
	// 系统级错误 (10000-19999)
	ErrInvalidParams    = 10001 // 无效的参数
	ErrInvalidJSON      = 10002 // 无效的JSON格式
	ErrMethodNotAllowed = 10003 // 方法不允许

	// 认证错误 (20000-29999)
	ErrUnauthorized      = 20001 // 未授权
	ErrInvalidToken      = 20002 // 无效的令牌
	ErrTokenExpired      = 20003 // 令牌过期
	ErrInvalidCredential = 20004 // 无效的凭证

	// 权限错误 (30000-39999)
	ErrPermissionDenied = 30001 // 权限不足
	ErrAccessForbidden  = 30002 // 访问被禁止

	// 资源错误 (40000-49999)
	ErrResourceNotFound = 40001 // 资源不存在
	ErrNodeNotFound     = 40002 // 节点不存在
	ErrLogNotFound      = 40003 // 日志不存在

	// 服务器错误 (50000-59999)
	ErrInternalServer     = 50001 // 服务器内部错误
	ErrDatabaseOperation  = 50002 // 数据库操作错误
	ErrNetworkError      = 50003 // 网络错误
	ErrWebSocketError    = 50004 // WebSocket错误
)

// 错误信息映射
var errorMessages = map[int]string{
	ErrInvalidParams:     "无效的参数",
	ErrInvalidJSON:       "无效的JSON格式",
	ErrMethodNotAllowed:  "方法不允许",
	ErrUnauthorized:      "未授权",
	ErrInvalidToken:      "无效的令牌",
	ErrTokenExpired:      "令牌过期",
	ErrInvalidCredential: "无效的凭证",
	ErrPermissionDenied:  "权限不足",
	ErrAccessForbidden:   "访问被禁止",
	ErrResourceNotFound:  "资源不存在",
	ErrNodeNotFound:      "节点不存在",
	ErrLogNotFound:       "日志不存在",
	ErrInternalServer:    "服务器内部错误",
	ErrDatabaseOperation: "数据库操作错误",
	ErrNetworkError:     "网络错误",
	ErrWebSocketError:   "WebSocket错误",
}

// GetErrorMessage 获取错误信息
func GetErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
