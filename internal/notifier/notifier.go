package notifier

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"monitor/internal/config"
	"monitor/internal/monitor"
)

// Manager 通知管理器（支持多种通知渠道）
type Manager struct {
	notifiers    []Notifier
	stateTracker *StateTracker
	config       *config.NotifierConfig
	mu           sync.RWMutex
}

// NewManager 创建通知管理器
func NewManager(cfg *config.NotifierConfig) (*Manager, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("通知功能未启用")
	}

	manager := &Manager{
		notifiers:    make([]Notifier, 0),
		stateTracker: NewStateTracker(cfg),
		config:       cfg,
	}

	// 初始化企业微信通知器
	if cfg.WeCom.Enabled {
		if cfg.WeCom.WebhookURL == "" {
			return nil, fmt.Errorf("企业微信 webhook_url 不能为空")
		}
		wecom, err := NewWeComNotifier(&cfg.WeCom)
		if err != nil {
			return nil, fmt.Errorf("初始化企业微信通知器失败: %w", err)
		}
		manager.notifiers = append(manager.notifiers, wecom)
		log.Printf("[Notifier] 企业微信通知器已启用")
	}

	if len(manager.notifiers) == 0 {
		return nil, fmt.Errorf("未配置任何通知渠道")
	}

	return manager, nil
}

// NotifyIfNeeded 检查是否需要发送告警，如需要则异步发送
func (m *Manager) NotifyIfNeeded(ctx context.Context, result *monitor.ProbeResult) {
	if m == nil {
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// 使用 StateTracker 判断是否需要告警
	alert := m.stateTracker.CheckAndBuildAlert(result)
	if alert == nil {
		return // 无需告警
	}

	// 异步发送通知（不阻塞探测流程）
	for _, notifier := range m.notifiers {
		n := notifier // 避免闭包问题
		go func() {
			if err := n.Send(ctx, alert); err != nil {
				log.Printf("[Notifier] 发送告警失败 %s-%s-%s: %v",
					alert.Provider, alert.Service, alert.Channel, err)
			} else {
				log.Printf("[Notifier] 告警已发送 %s-%s-%s: %s",
					alert.Provider, alert.Service, alert.Channel, alert.AlertType)
			}
		}()
	}
}

// UpdateConfig 热更新配置
func (m *Manager) UpdateConfig(cfg *config.NotifierConfig) {
	if m == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = cfg
	m.stateTracker.UpdateConfig(cfg)

	// TODO: 支持热更新通知器列表（企业微信 Webhook URL 等）
	log.Printf("[Notifier] 配置已更新")
}

// Close 关闭所有通知器
func (m *Manager) Close() error {
	if m == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, notifier := range m.notifiers {
		if err := notifier.Close(); err != nil {
			log.Printf("[Notifier] 关闭通知器失败: %v", err)
		}
	}

	return nil
}

// StateTracker 状态追踪器（追踪每个服务的状态变化）
type StateTracker struct {
	states map[string]*ServiceState // key: "provider-service-channel"
	mu     sync.RWMutex
	config *config.NotifierConfig
}

// ServiceState 服务状态
type ServiceState struct {
	LastStatus     int       // 上次状态（0/1/2）
	LastNotifyTime time.Time // 上次发送告警的时间
	FailureCount   int       // 连续失败次数
	FirstFailTime  time.Time // 首次失败时间
}

// NewStateTracker 创建状态追踪器
func NewStateTracker(cfg *config.NotifierConfig) *StateTracker {
	return &StateTracker{
		states: make(map[string]*ServiceState),
		config: cfg,
	}
}

// CheckAndBuildAlert 检查是否需要告警，并构造 Alert 对象
func (st *StateTracker) CheckAndBuildAlert(result *monitor.ProbeResult) *Alert {
	key := st.buildKey(result.Provider, result.Service, result.Channel)

	st.mu.Lock()
	defer st.mu.Unlock()

	// 获取或创建状态记录
	state, exists := st.states[key]
	if !exists {
		state = &ServiceState{
			LastStatus:     result.Status,
			LastNotifyTime: time.Time{}, // 零值，表示从未发送过告警
			FailureCount:   0,
		}
		st.states[key] = state

		// 首次探测，不发送告警
		if result.Status == StatusRed {
			state.FailureCount = 1
			state.FirstFailTime = time.Now()
		}
		return nil
	}

	// 更新状态前保存旧值
	previousStatus := state.LastStatus
	currentStatus := result.Status

	// 情况 1: 服务从正常变为不可用（绿/黄 → 红）
	if previousStatus != StatusRed && currentStatus == StatusRed {
		state.FailureCount = 1
		state.FirstFailTime = time.Now()
		state.LastStatus = currentStatus

		// 检查冷却期
		if st.isInCooldown(state) {
			return nil
		}

		// 构造告警
		alert := st.buildAlert(result, AlertTypeDown, previousStatus)
		state.LastNotifyTime = time.Now()
		return alert
	}

	// 情况 2: 服务持续不可用（红 → 红）
	if previousStatus == StatusRed && currentStatus == StatusRed {
		state.FailureCount++
		state.LastStatus = currentStatus

		// 检查是否达到持续不可用阈值
		if state.FailureCount == st.config.ContinuousFailureThreshold {
			// 检查冷却期
			if st.isInCooldown(state) {
				return nil
			}

			// 构造持续不可用告警
			alert := st.buildAlert(result, AlertTypeContinuousDown, previousStatus)
			alert.FailureCount = state.FailureCount
			state.LastNotifyTime = time.Now()
			return alert
		}

		// 未达到阈值，不告警
		return nil
	}

	// 情况 3: 服务从不可用恢复正常（红 → 绿/黄）
	if previousStatus == StatusRed && currentStatus != StatusRed {
		state.FailureCount = 0
		state.FirstFailTime = time.Time{}
		state.LastStatus = currentStatus

		// 检查冷却期
		if st.isInCooldown(state) {
			return nil
		}

		// 构造恢复告警
		alert := st.buildAlert(result, AlertTypeUp, previousStatus)
		state.LastNotifyTime = time.Now()
		return alert
	}

	// 情况 4: 正常状态之间的变化（绿 ↔ 黄）
	state.LastStatus = currentStatus
	return nil // 不告警
}

// isInCooldown 检查是否在冷却期内
func (st *StateTracker) isInCooldown(state *ServiceState) bool {
	if state.LastNotifyTime.IsZero() {
		return false // 从未发送过告警，不在冷却期
	}

	elapsed := time.Since(state.LastNotifyTime)
	return elapsed < st.config.MinNotifyIntervalDuration
}

// buildAlert 构造 Alert 对象
func (st *StateTracker) buildAlert(result *monitor.ProbeResult, alertType string, previousStatus int) *Alert {
	return &Alert{
		Provider:       result.Provider,
		Service:        result.Service,
		Channel:        result.Channel,
		Status:         result.Status,
		PreviousStatus: previousStatus,
		SubStatus:      string(result.SubStatus),
		Latency:        result.Latency,
		Timestamp:      result.Timestamp,
		AlertType:      alertType,
		FailureCount:   0, // 仅 continuous_down 时会设置
	}
}

// buildKey 生成状态映射的 key
func (st *StateTracker) buildKey(provider, service, channel string) string {
	return fmt.Sprintf("%s-%s-%s", provider, service, channel)
}

// UpdateConfig 更新配置
func (st *StateTracker) UpdateConfig(cfg *config.NotifierConfig) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.config = cfg
}

// GetState 获取服务状态（用于调试）
func (st *StateTracker) GetState(provider, service, channel string) *ServiceState {
	key := st.buildKey(provider, service, channel)
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.states[key]
}

// ClearOldStates 清理长时间未更新的状态（可选，防止内存泄漏）
func (st *StateTracker) ClearOldStates(maxAge time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()

	now := time.Now()
	for key, state := range st.states {
		// 如果最后通知时间超过 maxAge，且服务不在失败状态，则清理
		if !state.LastNotifyTime.IsZero() && now.Sub(state.LastNotifyTime) > maxAge && state.LastStatus != StatusRed {
			delete(st.states, key)
		}
	}
}
