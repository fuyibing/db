// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/fuyibing/log/v2"
	"github.com/fuyibing/log/v2/interfaces"
	"github.com/kataras/iris/v12"
	"github.com/teivah/onecontext"
	"xorm.io/xorm"
)

// Transaction handler.
type TransactionHandler func(ctx interface{}, sess *xorm.Session) error

// Bind context.
func Context(sess *xorm.Session, x interface{}) {
	// 1 get xorm session ctx
	v := reflect.ValueOf(sess).Elem()
	sessCtxRf := getUnexportedField(v.FieldByName("ctx"))
	sessCtx := sessCtxRf.(context.Context)

	// 2.1 return if nil.
	if x == nil {
		// base tracing ctx
		baseTracingCtx := log.NewContext()

		// merge tracing ctx wih xorm session ctx
		sessCtx, _ = onecontext.Merge(sessCtx, baseTracingCtx)
		sess.Context(sessCtx)
		return
	}
	// 2.2 context.Context.
	if c, ok := x.(context.Context); ok && c != nil {
		// merge context.Context wih xorm session ctx
		sessCtx, _ = onecontext.Merge(sessCtx, c)
		sess.Context(sessCtx)
		return
	}
	// 2.3 iris.Context.
	if c, ok := x.(iris.Context); ok && c != nil {
		if g := c.Values().Get(interfaces.OpenTracingKey); g != nil {
			// tracing ctx
			traceCtx := context.WithValue(context.TODO(), interfaces.OpenTracingKey, g.(interfaces.TraceInterface))

			// merge tracing ctx wih xorm session ctx
			sessCtx, _ = onecontext.Merge(sessCtx, traceCtx)
			sess.Context(sessCtx)
		}
		return
	}
}

// Return master connection.
func Master() *xorm.Session {
	engine := Config.engines.Master()
	engine.EnableSessionID(true)
	return engine.NewSession()
}

// Return master connection and bind context.
func MasterContext(ctx interface{}) *xorm.Session {
	sess := Master()
	Context(sess, ctx)
	return sess
}

// Return slave connection.
func Slave() *xorm.Session {
	if Config.slaveEnable {
		engine := Config.engines.Slave()
		engine.EnableSessionID(true)
		return engine.NewSession()
	}
	return Master()
}

// Return slave connection and bind context.
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
// Rollback when error return by handler, all handler executed with liner.
func TransactionWithSession(ctx interface{}, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	// 1. Select master connection if session not specified.
	if sess == nil {
		sess = MasterContext(ctx)
	}
	// 2. Begin transaction.
	if err = sess.Begin(); err != nil {
		return
	}
	// 3. Defer operation.
	defer func() {
		// 3.1 Catch panic.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			log.Errorfc(ctx, "[SQL] transaction panic: %s.", err.Error())
		}
		// 3.2 End transaction.
		if err != nil {
			// 3.2.1 rollback.
			_ = sess.Rollback()
		} else {
			// 3.2.2 commit.
			_ = sess.Commit()
		}
	}()
	// 4. call handler by liner.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}
	return
}

// get private unexported field
func getUnexportedField(field reflect.Value) interface{} {
	rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	return rf.Interface()
}
