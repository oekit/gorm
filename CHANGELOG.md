# Changelog — oekit/gorm

由 ekit v0.34.0（git.enjoye.top/enjoydream/ekit，冻结维护）拆分而来。所有 notable 变更记录如下。

## [v0.2.0] — 2026-10-15

可观测性统一轮。

### Breaking
- `gormx.Open` **不再默认注入 OTel 追踪插件**，go.mod 不再依赖 go.opentelemetry.io/*
- 删除 `WithoutTracing()` 选项（默认已不注入，选项失去意义）
- 删除 `NewOtelPlugin` / `OtelOption` 及全套 Otel* 选项（已迁至
  `github.com/oekit/tracing/gormotel`：`gormotel.NewPlugin` / `gormotel.Option`）

### 迁移
- 原来传 `WithoutTracing()` 的调用：直接删掉该参数
- 需要 SQL Span 的服务：`db.Use(gormotel.NewPlugin())` 显式挂载（无 SDK 时 noop）

## [v0.1.8] — 2026-10-15

- v0.1.8

## [v0.1.7] — 2026-10-15

- v0.1.7

## [v0.1.6] — 2026-10-15

- v0.1.6

## [v0.1.5] — 2026-10-15

- v0.1.5

## [v0.1.4] — 2026-10-15

- v0.1.4

## [v0.1.3] — 2026-10-15

- v0.1.3

## [v0.1.2] — 2026-10-15

- v0.1.2

## [v0.1.1] — 2026-10-15

- v0.1.1

## [v0.1.0] — 2026-10-15

- v0.1.0
