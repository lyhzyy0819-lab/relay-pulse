package config

// GetDefaultMessageTemplates 返回默认消息模板
func GetDefaultMessageTemplates() *MessageTemplates {
	return &MessageTemplates{
		Down: &MessageTemplate{
			Title:   "⚠️ 服务不可用告警",
			Content: defaultDownTemplate,
		},
		Up: &MessageTemplate{
			Title:   "✅ 服务恢复告警",
			Content: defaultUpTemplate,
		},
		ContinuousDown: &MessageTemplate{
			Title:   "🔴 服务持续不可用告警",
			Content: defaultContinuousDownTemplate,
		},
	}
}

const defaultDownTemplate = `> **服务商**: {{.Provider}}
> **服务**: {{.Service}}
{{- if .Channel}}
> **通道**: {{.Channel}}
{{- end}}
> **当前状态**: {{.StatusEmoji}} {{.StatusName}}
{{- if .SubStatusName}}
> **失败原因**: {{.SubStatusName}} (HTTP {{.HTTPStatusHint}})
{{- end}}
> **告警时间**: {{.Timestamp}}

*来自 RelayPulse 监控*`

const defaultUpTemplate = `> **服务商**: {{.Provider}}
> **服务**: {{.Service}}
{{- if .Channel}}
> **通道**: {{.Channel}}
{{- end}}
> **当前状态**: {{.StatusEmoji}} {{.StatusName}}
{{- if gt .Latency 0}}
> **响应延迟**: {{.Latency}} ms
{{- end}}
> **恢复时间**: {{.Timestamp}}

*来自 RelayPulse 监控*`

const defaultContinuousDownTemplate = `> **服务商**: {{.Provider}}
> **服务**: {{.Service}}
{{- if .Channel}}
> **通道**: {{.Channel}}
{{- end}}
> **当前状态**: {{.StatusEmoji}} {{.StatusName}}
{{- if gt .FailureCount 0}}
> **连续失败**: {{.FailureCount}} 次
{{- end}}
{{- if .SubStatusName}}
> **失败原因**: {{.SubStatusName}} (HTTP {{.HTTPStatusHint}})
{{- end}}
> **告警时间**: {{.Timestamp}}

*来自 RelayPulse 监控*`
