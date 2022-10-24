// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"context"
	"fmt"
	"github.com/fuyibing/log/v3"
	"xorm.io/xorm"
)

// TransactionHandler
// 事务回调.
type TransactionHandler func(ctx context.Context, sess *xorm.Session) error

// Transaction
// 执行事务.
func Transaction(ctx context.Context, handlers ...TransactionHandler) error {
	sess := Connector.GetMasterWithContext(ctx)
	return TransactionWithSession(ctx, sess, handlers...)
}

// TransactionWithSession
// 执行事务.
func TransactionWithSession(ctx context.Context, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	// 1. 开启事务.
	if err = sess.Begin(); err != nil {
		return
	}

	// 2. 事务结束.
	defer func() {
		// 2.1 运行异常.
		if r := recover(); r != nil {
			log.Panicfc(ctx, "panic on database transaction: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", err)
			}
		}

		// 2.2 结束事务.
		if err == nil {
			_ = sess.Commit()
		} else {
			_ = sess.Rollback()
		}
	}()

	// 3. 顺序执行.
	//    任一结果返回 error 时, 退出执行.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}
	return
}
