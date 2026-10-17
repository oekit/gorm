package gormx

import (
	"fmt"
	"time"

	"github.com/oekit/logx"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var gormLog = logx.Delegating("gormx")

// OpenOption 配置 Open 的行为（函数选项模式）。
type OpenOption func(*openConfig)

type openConfig struct {
	gormOptions   []gorm.Option
	maxIdleConns  int
	maxOpenConns  int
	connLifetime  time.Duration
	connIdleTime  time.Duration
	configurePool bool
}

// WithGormOptions 追加底层 gorm.Option。
func WithGormOptions(opts ...gorm.Option) OpenOption {
	return func(c *openConfig) {
		c.gormOptions = append(c.gormOptions, opts...)
	}
}

// WithConnectionPool 配置连接池参数（避免 OpenWithPool 的位置参数易错问题）。
func WithConnectionPool(maxIdleConns, maxOpenConns int, connMaxLifetime time.Duration) OpenOption {
	return func(c *openConfig) {
		c.maxIdleConns = maxIdleConns
		c.maxOpenConns = maxOpenConns
		c.connLifetime = connMaxLifetime
		c.configurePool = true
	}
}

// WithConnMaxIdleTime 配置连接最大空闲时间。
// 空闲超过此时间的连接会被关闭，防止 stale 连接积累。
func WithConnMaxIdleTime(d time.Duration) OpenOption {
	return func(c *openConfig) {
		c.connIdleTime = d
		c.configurePool = true
	}
}

// Open 使用指定的 dialector 打开数据库连接，自动注入默认日志。
//
// 默认行为：
//   - 注入默认 Logger（Info 级别）
//   - 零 tracing 注入：需要 SQL Span 经 github.com/oekit/tracing/gormotel
//     显式挂载
//
// 可通过选项配置连接池等：
//
//	db, err := gormx.Open(
//	    mysql.Open(dsn),
//	    gormx.WithConnectionPool(10, 100, time.Hour),
//	)
//
// 需要 OTel 追踪时（装饰器显式接入，无 SDK 时 noop）：
//
//	db, err := gormx.Open(mysql.Open(dsn))
//	if err := db.Use(gormotel.NewPlugin()); err != nil { ... }
func Open(dialector gorm.Dialector, opts ...OpenOption) (*gorm.DB, error) {
	cfg := &openConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	defaultLogger := NewDefaultLogger()

	gormCfg := gorm.Config{
		Logger: defaultLogger.LogMode(gormLogger.Info),
	}
	for _, gormOpt := range cfg.gormOptions {
		if err := gormOpt.Apply(&gormCfg); err != nil {
			return nil, fmt.Errorf("apply gorm option: %w", err)
		}
	}

	db, err := gorm.Open(dialector, &gormCfg)
	if err != nil {
		gormLog.Error("failed opening connection to db", "err", err)
		return nil, err
	}

	if cfg.configurePool {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
		sqlDB.SetConnMaxLifetime(cfg.connLifetime)
		if cfg.connIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(cfg.connIdleTime)
		}
	}

	return db, nil
}
