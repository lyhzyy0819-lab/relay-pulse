/**
 * 根据可用率计算渐变颜色
 * - 60%以下 → 灰色
 * - 60%-80% → 灰到浅橙渐变
 * - 80%-100% → 浅橙到深橙渐变
 * - -1（无数据）→ 灰色
 */

import type { CSSProperties } from 'react';

// 颜色常量
const GRAY = { r: 107, g: 114, b: 128 };        // #6b7280（不可用）
const LIGHT_ORANGE = { r: 251, g: 146, b: 60 }; // #fb923c（降级）
const DEEP_ORANGE = { r: 234, g: 88, b: 12 };   // #ea580c（可用）
const NO_DATA_GRAY = { r: 148, g: 163, b: 184 }; // #94a3b8（无数据）

interface RGB {
  r: number;
  g: number;
  b: number;
}

/**
 * 线性插值两个颜色
 */
function lerpColor(color1: RGB, color2: RGB, t: number): string {
  const r = Math.round(color1.r + (color2.r - color1.r) * t);
  const g = Math.round(color1.g + (color2.g - color1.g) * t);
  const b = Math.round(color1.b + (color2.b - color1.b) * t);
  return `rgb(${r}, ${g}, ${b})`;
}

/**
 * 根据可用率返回背景颜色（CSS color string）
 */
export function availabilityToColor(availability: number): string {
  // 无数据
  if (availability < 0) {
    return `rgb(${NO_DATA_GRAY.r}, ${NO_DATA_GRAY.g}, ${NO_DATA_GRAY.b})`;
  }

  // 60%以下 → 灰色
  if (availability < 60) {
    return `rgb(${GRAY.r}, ${GRAY.g}, ${GRAY.b})`;
  }

  // 60%-80% → 灰到浅橙渐变
  if (availability < 80) {
    const t = (availability - 60) / 20;
    return lerpColor(GRAY, LIGHT_ORANGE, t);
  }

  // 80%-100% → 浅橙到深橙渐变
  const t = (availability - 80) / 20;
  return lerpColor(LIGHT_ORANGE, DEEP_ORANGE, t);
}

/**
 * 根据可用率返回 Tailwind 兼容的 style 对象
 */
export function availabilityToStyle(availability: number): CSSProperties {
  return {
    backgroundColor: availabilityToColor(availability),
  };
}

/**
 * 根据延迟计算渐变颜色
 * - 延迟越低越好（与可用率相反）
 * - 基于 slow_latency 阈值进行相对渐变
 *
 * 渐变逻辑：
 * - latency <= 0 → 灰色（无数据）
 * - latency < 30% 阈值 → 深橙（优秀）
 * - 30%-100% 阈值 → 深橙到浅橙渐变（良好）
 * - 100%-200% 阈值 → 浅橙到灰渐变（较慢）
 * - >= 200% 阈值 → 灰色（很慢）
 */
export function latencyToColor(latency: number, slowLatencyMs: number): string {
  // 无数据或配置无效
  if (latency <= 0 || slowLatencyMs <= 0) {
    return `rgb(${NO_DATA_GRAY.r}, ${NO_DATA_GRAY.g}, ${NO_DATA_GRAY.b})`;
  }

  const ratio = latency / slowLatencyMs;

  // < 30% 阈值 → 深橙
  if (ratio < 0.3) {
    return `rgb(${DEEP_ORANGE.r}, ${DEEP_ORANGE.g}, ${DEEP_ORANGE.b})`;
  }

  // 30%-100% 阈值 → 深橙到浅橙渐变
  if (ratio < 1) {
    const t = (ratio - 0.3) / 0.7;
    return lerpColor(DEEP_ORANGE, LIGHT_ORANGE, t);
  }

  // 100%-200% 阈值 → 浅橙到灰渐变
  if (ratio < 2) {
    const t = (ratio - 1) / 1;
    return lerpColor(LIGHT_ORANGE, GRAY, t);
  }

  // >= 200% 阈值 → 灰色
  return `rgb(${GRAY.r}, ${GRAY.g}, ${GRAY.b})`;
}
