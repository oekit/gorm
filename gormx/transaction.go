package gormx

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// TXFunc 是在事务中执行的函数类型，接收事务内的 *gorm.DB 并返回 error。
type TXFunc func(tx *gorm.DB) error

// InTransactionContext 在带有 context 的数据库事务中执行指定函数。
func InTransactionContext(ctx context.Context, db *gorm.DB, fn TXFunc) error {
	return db.WithContext(ctx).Transaction(fn)
}

// IsUniqueViolation 判断给定 error 是否为唯一键冲突。
//
// 覆盖三类来源：
//  1. gorm 翻译错误（TranslateError: true 时驱动错误映射为 gorm.ErrDuplicatedKey）
//  2. 驱动原始报文：MySQL「Duplicate entry」、PostgreSQL「duplicate key value」/「23505」、
//     SQLite「UNIQUE constraint failed: <table>.<col>」（mattn 与 glebarez 驱动同格式，
//     均不映射 gorm.ErrDuplicatedKey，只能按报文匹配）
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "duplicate key value") ||
		strings.Contains(msg, "23505") ||
		strings.Contains(msg, "UNIQUE constraint failed")
}
