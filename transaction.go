// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"context"
	"fmt"
	"github.com/fuyibing/log/v8"
	"xorm.io/xorm"
)

type TransactionHandler func(ctx context.Context, sess *xorm.Session) error

func Transaction(ctx context.Context, handlers ...TransactionHandler) error {
	sess := Connector.GetMasterWithContext(ctx)
	defer func() { _ = sess.Close() }()
	err := TransactionWithSession(ctx, sess, handlers...)
	return err
}

func TransactionWithSession(ctx context.Context, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	if err = sess.Begin(); err != nil {
		return
	}

	// Transaction result detect.
	defer func() {
		if r := recover(); r != nil {
			log.Fatalfc(ctx, "panic on database transaction: %v", r)

			if err == nil {
				err = fmt.Errorf("%v", err)
			}
		}

		if err == nil {
			_ = sess.Commit()
		} else {
			_ = sess.Rollback()
		}
	}()

	// Iterate registered handlers. Break if error returned
	// from any handler.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}
	return
}
