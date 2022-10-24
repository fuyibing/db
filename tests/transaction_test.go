// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package tests

import (
	"context"
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"testing"
	"xorm.io/xorm"
)

func TestTransaction(t *testing.T) {
	ctx := trace.New()
	log.Infofc(ctx, "testing.Transaction begin")

	err := db.Transaction(ctx, tx1, tx2, tx3)
	log.Infofc(ctx, "testing.Transaction: end %v", err)
}

func tx1(_ context.Context, sess *xorm.Session) (err error) {
	_, err = NewExampleService(sess).GetById(1)
	return
}

func tx2(_ context.Context, sess *xorm.Session) (err error) {
	_, err = NewExampleService(sess).GetById(2)
	return
}

func tx3(_ context.Context, sess *xorm.Session) (err error) {
	_, err = NewExampleService(sess).GetById(3)
	return
}
