// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"xorm.io/xorm"

	"github.com/fuyibing/log"
)

// Transaction handler.
type TransactionHandler func(ctx interface{}, sess *xorm.Session) error

// Set context on connection session.
func Context(sess *xorm.Session, x interface{}) {
	// return if nil.
	if x == nil {
		return
	}
	// context.Context.
	if c, ok := x.(context.Context); ok && c != nil {
		sess.Context(c)
		return
	}
	// iris.Context.
	if c, ok := x.(iris.Context); ok && c != nil {
		if g := c.Values().Get(log.OpenTracingContext); g != nil {
			sess.Context(g.(context.Context))
		}
		return
	}
}

// Return master connection session.
func Master() *xorm.Session {
	return Config.engines.Master().NewSession()
}

// Return master connection session with context.
func MasterContext(ctx interface{}) *xorm.Session {
	sess := Master()
	Context(sess, ctx)
	return sess
}

// Return slave connection session.
func Slave() *xorm.Session {
	return Config.engines.Slave().NewSession()
}

// Return slave connection session with context.
func SlaveContext(ctx interface{}) *xorm.Session {
	sess := Slave()
	Context(sess, ctx)
	return sess
}

// Run transaction.
func Transaction(ctx interface{}, handlers ...TransactionHandler) (err error) {
	return TransactionWithSession(ctx, nil, handlers...)
}

// Run transaction.
//
//   ctx := log.NewContext()
//   sess := xdb.MasterContext(tracing)
//   if err := xdb.TransactionWithSession(tracing, sess, func(ctx interface{}, sess *xorm.Session) error {
//       // logic
//   }, func(ctx interface{}, sess *xorm.Session) error {
//       // logic
//   }, func(ctx interface{}, sess *xorm.Session) error {
//       // logic
//   }); err != nil {
//       println("Transaction error - ", err.Error())
//   }
//
func TransactionWithSession(ctx interface{}, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	// Create master connection session if not specified.
	if sess == nil {
		sess = Master()
	}
	// Transaction begin error.
	if err = sess.Begin(); err != nil {
		return
	}
	// End transaction.
	defer func() {
		// Catch panic.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		// End transaction.
		if err != nil {
			// Rollback.
			_ = sess.Rollback()
		} else {
			// Commit.
			_ = sess.Commit()
		}
	}()
	// Loop handlers.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}
	// Completed.
	return
}
