package notifier

import (
	"strings"
	"testing"

	"monitor/internal/config"
)

func TestMessageBuilder_BuildMessage_DefaultTemplate(t *testing.T) {
	// 使用默认模板
	templates := config.GetDefaultMessageTemplates()
	builder, err := NewMessageBuilder(templates)
	if err != nil {
		t.Fatalf("创建 MessageBuilder 失败: %v", err)
	}

	// 测试用例 1: down 告警
	t.Run("down alert", func(t *testing.T) {
		alert := &Alert{
			Provider:  "Code-CLI",
			Service:   "cc",
			Channel:   "vip-channel",
			Status:    StatusRed,
			SubStatus: "rate_limit",
			Timestamp: 1735559123,
			AlertType: AlertTypeDown,
		}

		msg, err := builder.BuildMessage(alert)
		if err != nil {
			t.Fatalf("构造消息失败: %v", err)
		}

		// 验证消息包含关键信息
		if !strings.Contains(msg, "Code-CLI") {
			t.Errorf("消息不包含服务商名称")
		}
		if !strings.Contains(msg, "vip-channel") {
			t.Errorf("消息不包含通道信息")
		}
		if !strings.Contains(msg, "⚠️ 服务不可用告警") {
			t.Errorf("消息不包含标题")
		}
		if !strings.Contains(msg, "限流") {
			t.Errorf("消息不包含失败原因")
		}
	})

	// 测试用例 2: up 告警
	t.Run("up alert", func(t *testing.T) {
		alert := &Alert{
			Provider:  "Code-CLI",
			Service:   "cc",
			Status:    StatusGreen,
			Latency:   234,
			Timestamp: 1735559123,
			AlertType: AlertTypeUp,
		}

		msg, err := builder.BuildMessage(alert)
		if err != nil {
			t.Fatalf("构造消息失败: %v", err)
		}

		if !strings.Contains(msg, "✅ 服务恢复告警") {
			t.Errorf("消息不包含标题")
		}
		if !strings.Contains(msg, "234 ms") {
			t.Errorf("消息不包含延迟信息")
		}
	})

	// 测试用例 3: continuous_down 告警
	t.Run("continuous_down alert", func(t *testing.T) {
		alert := &Alert{
			Provider:     "Code-CLI",
			Service:      "cc",
			Status:       StatusRed,
			SubStatus:    "server_error",
			FailureCount: 5,
			Timestamp:    1735559123,
			AlertType:    AlertTypeContinuousDown,
		}

		msg, err := builder.BuildMessage(alert)
		if err != nil {
			t.Fatalf("构造消息失败: %v", err)
		}

		if !strings.Contains(msg, "🔴 服务持续不可用告警") {
			t.Errorf("消息不包含标题")
		}
		if !strings.Contains(msg, "5 次") {
			t.Errorf("消息不包含连续失败次数")
		}
	})

	// 测试用例 4: 无 Channel 的告警（条件判断测试）
	t.Run("alert without channel", func(t *testing.T) {
		alert := &Alert{
			Provider:  "Code-CLI",
			Service:   "cc",
			Channel:   "", // 空 Channel
			Status:    StatusRed,
			Timestamp: 1735559123,
			AlertType: AlertTypeDown,
		}

		msg, err := builder.BuildMessage(alert)
		if err != nil {
			t.Fatalf("构造消息失败: %v", err)
		}

		// 验证消息不包含 "通道" 字段
		if strings.Contains(msg, "**通道**") {
			t.Errorf("消息不应包含通道字段")
		}
	})
}

func TestMessageBuilder_CustomTemplate(t *testing.T) {
	// 自定义模板
	customTemplates := &config.MessageTemplates{
		Down: &config.MessageTemplate{
			Title:   "🚨 紧急告警",
			Content: "{{.Provider}} {{.Service}} 故障",
		},
		Up: &config.MessageTemplate{
			Title:   "✅ 恢复正常",
			Content: "{{.Provider}} {{.Service}} 已恢复",
		},
		ContinuousDown: &config.MessageTemplate{
			Title:   "🔥 严重告警",
			Content: "{{.Provider}} {{.Service}} 连续失败 {{.FailureCount}} 次",
		},
	}

	builder, err := NewMessageBuilder(customTemplates)
	if err != nil {
		t.Fatalf("创建 MessageBuilder 失败: %v", err)
	}

	alert := &Alert{
		Provider:  "TestProvider",
		Service:   "TestService",
		Status:    StatusRed,
		Timestamp: 1735559123,
		AlertType: AlertTypeDown,
	}

	msg, err := builder.BuildMessage(alert)
	if err != nil {
		t.Fatalf("构造消息失败: %v", err)
	}

	if !strings.Contains(msg, "🚨 紧急告警") {
		t.Errorf("消息不包含自定义标题")
	}
	if !strings.Contains(msg, "TestProvider TestService 故障") {
		t.Errorf("消息格式不正确")
	}
}

func TestMessageBuilder_InvalidTemplate(t *testing.T) {
	// 无效模板（语法错误）
	invalidTemplates := &config.MessageTemplates{
		Down: &config.MessageTemplate{
			Title:   "测试",
			Content: "{{.Provider {{.Service}}", // 语法错误
		},
		Up: &config.MessageTemplate{
			Title:   "测试",
			Content: "正常",
		},
		ContinuousDown: &config.MessageTemplate{
			Title:   "测试",
			Content: "正常",
		},
	}

	_, err := NewMessageBuilder(invalidTemplates)
	if err == nil {
		t.Errorf("应该返回语法错误")
	}
}
