package notifier

import "context"

// Notifier 通知器接口（支持多种通知渠道：企业微信、邮件、Slack 等）
type Notifier interface {
	// Send 发送告警通知
	Send(ctx context.Context, alert *Alert) error

	// Close 关闭通知器，清理资源
	Close() error
}

// Alert 告警结构
type Alert struct {
	// 服务标识
	Provider string // 服务商（如 "Code-CLI"）
	Service  string // 服务类型（如 "cc"）
	Channel  string // 业务通道（如 "vip-channel"）

	// 状态信息
	Status         int    // 当前状态（0=红色不可用, 1=绿色正常, 2=黄色降级）
	PreviousStatus int    // 上次状态
	SubStatus      string // 细分状态（rate_limit、server_error、network_error 等）

	// 性能指标
	Latency int // 响应延迟（毫秒）

	// 告警元信息
	Timestamp    int64  // 告警时间（Unix 时间戳）
	AlertType    string // 告警类型："down"（服务不可用）、"up"（服务恢复）、"continuous_down"（持续不可用）
	FailureCount int    // 连续失败次数（仅 continuous_down 时有意义）
}

// AlertType 常量
const (
	AlertTypeDown           = "down"            // 服务从正常变为不可用
	AlertTypeUp             = "up"              // 服务从不可用恢复正常
	AlertTypeContinuousDown = "continuous_down" // 服务持续不可用超过阈值
)

// Status 常量
const (
	StatusRed    = 0 // 红色：不可用
	StatusGreen  = 1 // 绿色：正常
	StatusYellow = 2 // 黄色：降级（延迟过高）
)

// StatusName 返回状态名称
func StatusName(status int) string {
	switch status {
	case StatusGreen:
		return "正常"
	case StatusYellow:
		return "降级"
	case StatusRed:
		return "不可用"
	default:
		return "未知"
	}
}

// StatusEmoji 返回状态对应的 emoji
func StatusEmoji(status int) string {
	switch status {
	case StatusGreen:
		return "🟢"
	case StatusYellow:
		return "🟡"
	case StatusRed:
		return "🔴"
	default:
		return "⚪"
	}
}

// SubStatusName 返回细分状态的中文名称
func SubStatusName(subStatus string) string {
	switch subStatus {
	case "rate_limit":
		return "限流"
	case "server_error":
		return "服务器错误"
	case "client_error":
		return "客户端错误"
	case "auth_error":
		return "认证失败"
	case "invalid_request":
		return "请求参数错误"
	case "network_error":
		return "网络错误"
	case "content_mismatch":
		return "内容校验失败"
	case "slow_latency":
		return "响应慢"
	default:
		return ""
	}
}
