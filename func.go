// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/fuyibing/log/v2"
	"github.com/fuyibing/log/v2/interfaces"
	"github.com/kataras/iris/v12"
	"xorm.io/xorm"
)

// 事务回调.
type TransactionHandler func(ctx interface{}, sess *xorm.Session) error

// 绑定Context.
func Context(sess *xorm.Session, x interface{}) {
	// return if nil.
	if x == nil {
		sess.Context(log.NewContext())
		return
	}
	// context.Context.
	if c, ok := x.(context.Context); ok && c != nil {
		sess.Context(c)
		return
	}
	// iris.Context.
	if c, ok := x.(iris.Context); ok && c != nil {
		if g := c.Values().Get(interfaces.OpenTracingKey); g != nil {
			sess.Context(context.WithValue(context.TODO(), interfaces.OpenTracingKey, g.(interfaces.TraceInterface)))
		}
		return
	}
}

// 读取主库连结.
func Master() *xorm.Session {
	return Config.engines.Master().NewSession()
}

// 读取主库连结, 并绑定Context.
func MasterContext(ctx interface{}) *xorm.Session {
	sess := Master()
	Context(sess, ctx)
	return sess
}

// 读取从库连结.
func Slave() *xorm.Session {
	if Config.slaveEnable {
		return Config.engines.Slave().NewSession()
	}
	return Master()
}

// 读取从库连结, 并绑定Context.
func SlaveContext(ctx interface{}) *xorm.Session {
	sess := Slave()
	Context(sess, ctx)
	return sess
}

// 执行事务.
func Transaction(ctx interface{}, handlers ...TransactionHandler) (err error) {
	return TransactionWithSession(ctx, nil, handlers...)
}

// 执行事务.
//
// 在执行事务过程中, 所有操作都在同一个主协程序中串行处理, 当任意一个回调返回error时, 将
// 触发Rollback回滚.
func TransactionWithSession(ctx interface{}, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
	// 1. 选择主库.
	// 仅在未指定连结时生效.
	if sess == nil {
		sess = Master()
	}
	// 2. 开始事务.
	if err = sess.Begin(); err != nil {
		return
	}
	// 3. 捕获Panic.
	// 当脚本在运行过程中, 有可能产生运行Panic, 此处理捕获panic并
	// 触发回滚.
	defer func() {
		// 3.1 捕获Panic.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
			log.Errorfc(ctx, "[SQL] transaction panic: %s.", err.Error())
		}
		// 3.2 结束事务.
		if err != nil {
			// 3.2.1 回滚.
			_ = sess.Rollback()
		} else {
			// 3.2.2 提交.
			_ = sess.Commit()
		}
	}()
	// 4. 串行执行回调.
	for _, handler := range handlers {
		if err = handler(ctx, sess); err != nil {
			break
		}
	}
	return
}
