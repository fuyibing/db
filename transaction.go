// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"context"
	"fmt"
	"xorm.io/xorm"

	"github.com/fuyibing/log/v5"
)

type (
	// TransactionHandler
	// 事务Handler回调.
	TransactionHandler func(ctx context.Context, sess *xorm.Session) error
)

// Transaction
// 创建事务.
func Transaction(ctx context.Context, handlers ...TransactionHandler) error {
	sess := Connector.GetMasterWithContext(ctx)

	defer func() {
		_ = sess.Close()
	}()

	return TransactionWithSession(ctx, sess, handlers...)
}

// TransactionWithSession
// 基于指定连接创建事务.
func TransactionWithSession(ctx context.Context, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	// 开启事务.
	if err = sess.Begin(); err != nil {
		return
	}

	// 事务完成.
	defer func() {
		// 捕获异常.
		if r := recover(); r != nil {
			if spa, exists := log.Span(ctx); exists {
				spa.Logger().Fatal("transaction fatal: %v", r)
			} else {
				log.Fatal("transaction fatal: %v", r)
			}

			// 覆盖错误.
			if err == nil {
				err = fmt.Errorf("%v", err)
			}
		}

		// 事务提交.
		if err == nil {
			_ = sess.Commit()
		} else {
			_ = sess.Rollback()
		}
	}()

	// 遍历事务回调.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}

	return
}
