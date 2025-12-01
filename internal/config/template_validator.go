package config

import (
	"fmt"
	"text/template"
)

// validateMessageTemplates 验证消息模板语法
func validateMessageTemplates(templates *MessageTemplates) error {
	if templates == nil {
		return nil
	}

	// 验证 down 模板
	if err := validateTemplate(templates.Down, "down"); err != nil {
		return err
	}

	// 验证 up 模板
	if err := validateTemplate(templates.Up, "up"); err != nil {
		return err
	}

	// 验证 continuous_down 模板
	if err := validateTemplate(templates.ContinuousDown, "continuous_down"); err != nil {
		return err
	}

	return nil
}

// validateTemplate 验证单个模板
func validateTemplate(tmpl *MessageTemplate, name string) error {
	if tmpl == nil {
		return nil
	}

	if tmpl.Content == "" {
		return fmt.Errorf("模板 %s 的 content 不能为空", name)
	}

	// 尝试编译模板
	_, err := template.New(name).Parse(tmpl.Content)
	if err != nil {
		return fmt.Errorf("模板 %s 语法错误: %w", name, err)
	}

	return nil
}
