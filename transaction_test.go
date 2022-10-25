// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package db

import (
	"context"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"testing"
	"xorm.io/xorm"
)

func ExampleTransaction() {
	// 1. 创建上下文.
	ctx := trace.New()
	log.Infofc(ctx, "example begin")

	// 2. 定义事务.
	//
	//    参数 tx1, tx2, tx3, tx4 遵循 TransactionHandler, 执行时
	//    按顺序执行. 当任一返回 error 时, 跳过未执行的 handler.
	err := Transaction(ctx, tx1, tx2, tx3, tx4)

	// 3. 事务出错.
	if err != nil {
		log.Infofc(ctx, "example error: %v", err)
		return
	}

	// 4. 成功完成.
	log.Infofc(ctx, "example completed")
}

func ExampleTransactionWithSession() {
	// 1. 创建上下文.
	ctx := trace.New()
	log.Infofc(ctx, "example begin")

	// 2. 获取连接
	//    并在结束时释放回池.
	sess := Connector.GetMasterWithContext(ctx)
	defer func() { _ = sess.Close() }()

	// 3. 定义事务.
	//
	//    参数 tx1, tx2, tx3, tx4 遵循 TransactionHandler, 执行时
	//    按顺序执行. 当任一返回 error 时, 跳过未执行的 handler.
	err := TransactionWithSession(ctx, sess, tx1, tx2, tx3, tx4)

	// 4. 事务出错.
	if err != nil {
		log.Infofc(ctx, "example error: %v", err)
		return
	}

	// 5. 成功完成.
	log.Infofc(ctx, "example completed")
}

func TestTransaction(t *testing.T) { ExampleTransaction() }

func TestTransactionWithSession(t *testing.T) { ExampleTransactionWithSession() }

func tx1(_ context.Context, sess *xorm.Session) (err error) {
	_, err = newExampleService(sess).GetById(1)
	return
}

func tx2(_ context.Context, sess *xorm.Session) (err error) {
	_, err = newExampleService(sess).GetById(1)
	return
}

func tx3(_ context.Context, sess *xorm.Session) (err error) {
	_, err = newExampleService(sess).GetById(1)
	return
}

func tx4(_ context.Context, sess *xorm.Session) (err error) {
	_, err = newExampleService(sess).GetById(1)
	return
}
