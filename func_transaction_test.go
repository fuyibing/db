// author: wsfuyibing <websearch@163.com>
// date: 2022-06-07

package db

import (
    "testing"

    "github.com/fuyibing/log/v2"
    "xorm.io/xorm"
)

func TestTransaction(t *testing.T) {
    ctx := log.NewContext()
    log.Infofc(ctx, "testing.TestTransaction")

    if err := Transaction(ctx, transaction1, transaction2); err != nil {
        t.Errorf("testing: %v.", err)
        return
    }
    t.Logf("testing.completed.")
}

func TestTransactionWithSession(t *testing.T) {
    ctx := log.NewContext()
    log.Infofc(ctx, "testing.TestTransaction")

    sess := MasterContext(ctx)
    if err := TransactionWithSession(ctx, sess, transaction1, transaction2); err != nil {
        t.Errorf("testing: %v.", err)
        return
    }
    t.Logf("testing.completed.")
}

func transaction1(ctx interface{}, sess *xorm.Session) error {
    _, err := NewExampleService(sess).GetById(1)
    return err
}

func transaction2(ctx interface{}, sess *xorm.Session) error {
    _, err := NewExampleService(sess).GetById(2)
    return err
}
