package notifier

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"monitor/internal/config"
)

// TemplateData 模板数据（暴露给 text/template）
type TemplateData struct {
	Provider       string
	Service        string
	Channel        string
	StatusName     string
	StatusEmoji    string
	SubStatusName  string
	HTTPStatusHint string
	Timestamp      string
	FailureCount   int
	Latency        int
}

// MessageBuilder 消息构造器
type MessageBuilder struct {
	templates     *config.MessageTemplates
	compiledCache map[string]*template.Template
}

// NewMessageBuilder 创建消息构造器
func NewMessageBuilder(templates *config.MessageTemplates) (*MessageBuilder, error) {
	if templates == nil {
		return nil, fmt.Errorf("templates 不能为 nil")
	}

	builder := &MessageBuilder{
		templates:     templates,
		compiledCache: make(map[string]*template.Template),
	}

	// 预编译所有模板
	if err := builder.compileTemplates(); err != nil {
		return nil, err
	}

	return builder, nil
}

// compileTemplates 编译所有模板
func (mb *MessageBuilder) compileTemplates() error {
	if err := mb.compileTemplate("down", mb.templates.Down.Content); err != nil {
		return err
	}

	if err := mb.compileTemplate("up", mb.templates.Up.Content); err != nil {
		return err
	}

	if err := mb.compileTemplate("continuous_down", mb.templates.ContinuousDown.Content); err != nil {
		return err
	}

	return nil
}

// compileTemplate 编译单个模板
func (mb *MessageBuilder) compileTemplate(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return fmt.Errorf("编译模板 %s 失败: %w", name, err)
	}
	mb.compiledCache[name] = tmpl
	return nil
}

// BuildMessage 构造 Markdown 消息
func (mb *MessageBuilder) BuildMessage(alert *Alert) (string, error) {
	// 准备模板数据
	data := mb.prepareTemplateData(alert)

	// 根据告警类型选择模板
	var tmpl *template.Template
	var title string

	switch alert.AlertType {
	case AlertTypeDown:
		tmpl = mb.compiledCache["down"]
		title = mb.templates.Down.Title
	case AlertTypeUp:
		tmpl = mb.compiledCache["up"]
		title = mb.templates.Up.Title
	case AlertTypeContinuousDown:
		tmpl = mb.compiledCache["continuous_down"]
		title = mb.templates.ContinuousDown.Title
	default:
		return "", fmt.Errorf("未知的告警类型: %s", alert.AlertType)
	}

	if tmpl == nil {
		return "", fmt.Errorf("模板未编译: %s", alert.AlertType)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	// 组装最终消息（标题 + 内容）
	finalMsg := fmt.Sprintf("## %s\n\n%s", title, buf.String())
	return finalMsg, nil
}

// prepareTemplateData 准备模板数据
func (mb *MessageBuilder) prepareTemplateData(alert *Alert) *TemplateData {
	timestamp := time.Unix(alert.Timestamp, 0).Format("2006-01-02 15:04:05")

	return &TemplateData{
		Provider:       alert.Provider,
		Service:        alert.Service,
		Channel:        alert.Channel,
		StatusName:     StatusName(alert.Status),
		StatusEmoji:    StatusEmoji(alert.Status),
		SubStatusName:  SubStatusName(alert.SubStatus),
		HTTPStatusHint: getHTTPStatusHint(alert.SubStatus),
		Timestamp:      timestamp,
		FailureCount:   alert.FailureCount,
		Latency:        alert.Latency,
	}
}

// getHTTPStatusHint 根据 SubStatus 返回 HTTP 状态码提示
func getHTTPStatusHint(subStatus string) string {
	switch subStatus {
	case "rate_limit":
		return "429"
	case "server_error":
		return "5xx"
	case "client_error":
		return "4xx"
	case "auth_error":
		return "401/403"
	case "invalid_request":
		return "400"
	case "network_error":
		return "网络错误"
	case "content_mismatch":
		return "2xx"
	default:
		return ""
	}
}
