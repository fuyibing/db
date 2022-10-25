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
//
// 参数 ctx 是日志上下文, 参数 sess 为DB连接. 当返回 error 后不再执行后续回调.
type TransactionHandler func(ctx context.Context, sess *xorm.Session) error

// Transaction
// 执行事务.
//
//   [13:40:26.000][ INFO][span-id=0.0] example begin
//   [13:40:26.026][ INFO][span-id=0.1] [SQL] BEGIN TRANSACTION
//   [13:40:26.031][ INFO][span-id=0.2] [SQL] SELECT * FROM `example` WHERE (`id` = ?) LIMIT 1, Args: [1]
//   [13:40:26.031][ERROR][span-id=0.3] Error 1146: Table 'schema.example' doesn't exist
//   [13:40:26.037][ INFO][span-id=0.4] [SQL] ROLLBACK
//   [13:40:26.037][ INFO][span-id=0.5] example error: Error 1146: Table 'schema.example' doesn't exist
func Transaction(ctx context.Context, handlers ...TransactionHandler) error {
	sess := Connector.GetMasterWithContext(ctx)
	defer func() { _ = sess.Close() }()
	err := TransactionWithSession(ctx, sess, handlers...)
	return err
}

// TransactionWithSession
// 执行事务.
//
//   [13:40:26.000][ INFO][span-id=0.0] example begin
//   [13:40:26.026][ INFO][span-id=0.1] [SQL] BEGIN TRANSACTION
//   [13:40:26.031][ INFO][span-id=0.2] [SQL] SELECT * FROM `example` WHERE (`id` = ?) LIMIT 1, Args: [1]
//   [13:40:26.031][ERROR][span-id=0.3] Error 1146: Table 'schema.example' doesn't exist
//   [13:40:26.037][ INFO][span-id=0.4] [SQL] ROLLBACK
//   [13:40:26.037][ INFO][span-id=0.5] example error: Error 1146: Table 'schema.example' doesn't exist
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
