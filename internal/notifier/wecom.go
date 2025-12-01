package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"monitor/internal/config"
)

// WeComNotifier 企业微信通知器
type WeComNotifier struct {
	webhookURL     string
	client         *http.Client
	retryCount     int
	messageBuilder *MessageBuilder
}

// NewWeComNotifier 创建企业微信通知器
func NewWeComNotifier(cfg *config.WeComConfig) (*WeComNotifier, error) {
	if cfg.WebhookURL == "" {
		return nil, fmt.Errorf("webhook_url 不能为空")
	}

	// 初始化消息构造器
	msgBuilder, err := NewMessageBuilder(cfg.Templates)
	if err != nil {
		return nil, fmt.Errorf("初始化消息构造器失败: %w", err)
	}

	return &WeComNotifier{
		webhookURL: cfg.WebhookURL,
		client: &http.Client{
			Timeout: cfg.TimeoutDuration,
		},
		retryCount:     cfg.RetryCount,
		messageBuilder: msgBuilder,
	}, nil
}

// Send 发送告警通知
func (w *WeComNotifier) Send(ctx context.Context, alert *Alert) error {
	// 使用新的消息构造器
	msg, err := w.messageBuilder.BuildMessage(alert)
	if err != nil {
		return fmt.Errorf("构造消息失败: %w", err)
	}

	// 构造请求体（企业微信 Webhook API 格式）
	reqBody := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": msg,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 发送 HTTP POST（支持重试）
	var lastErr error
	for i := 0; i <= w.retryCount; i++ {
		if err := w.post(ctx, bodyBytes); err == nil {
			return nil // 成功
		} else {
			lastErr = err
			if i < w.retryCount {
				// 指数退避重试
				sleepDuration := time.Second * time.Duration(i+1)
				log.Printf("[WeComNotifier] 发送失败，%v 后重试 (%d/%d): %v",
					sleepDuration, i+1, w.retryCount, err)
				time.Sleep(sleepDuration)
			}
		}
	}

	return fmt.Errorf("达到最大重试次数，最后错误: %w", lastErr)
}

// post 发送 HTTP POST 请求
func (w *WeComNotifier) post(ctx context.Context, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", w.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 企业微信 Webhook API 成功返回 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP 状态码异常: %d", resp.StatusCode)
	}

	// 解析响应（企业微信返回 {"errcode":0,"errmsg":"ok"}）
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("企业微信 API 返回错误: %s (errcode=%d)", result.ErrMsg, result.ErrCode)
	}

	return nil
}

// Close 关闭通知器
func (w *WeComNotifier) Close() error {
	// HTTP 客户端无需显式关闭
	return nil
}
