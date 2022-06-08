// author: wsfuyibing <websearch@163.com>
// date: 2022-06-07

package db

import (
    "fmt"

    "github.com/fuyibing/log/v2"
    "xorm.io/xorm"
)

// TransactionHandler
// 事务执行器.
//
//   // 注意事项
//   // 1. 任一执行器返回error时退出事务(忽略未执行的执行器).
//   // 2. 事务中不可使用协程(goroutine, 务必在同一个协程中完成全部执行器).
//
//   // 执行器 1.
//   func handler1(ctx interface{}, sess *xorm.Session) error {
//       // ...
//   }
//
//   // 执行器 2.
//   func handler2(ctx interface{}, sess *xorm.Session) error {
//       // ...
//   }
//
//   // 执行器 3.
//   func handler3(ctx interface{}, sess *xorm.Session) error {
//       // ...
//   }
//
//   // 事务用法.
//   func main(ctx context.Context){
//       if err := db.Transaction(ctx, handler1, handler2, handler3); err != nil {
//           // ...
//       }
//   }
type TransactionHandler func(ctx interface{}, sess *xorm.Session) error

// Transaction
// 事务打包.
func Transaction(ctx interface{}, handlers ...TransactionHandler) (err error) {
    return TransactionWithSession(ctx, nil, handlers...)
}

// TransactionWithSession
// 事务打包.
func TransactionWithSession(ctx interface{}, sess *xorm.Session, handlers ...TransactionHandler) (err error) {
    // 1. 获取连接.
    //    若未指定连接, 则获取默认DB实例的主库连接.
    if sess == nil {
        if ctx == nil {
            sess = Master()
        } else {
            sess = MasterContext(ctx)
        }
    }

    // 2. 打开事务.
    if log.Config.DebugOn() {
        log.Debugfc(ctx, "[transaction] begin.")
    }
    if err = sess.Begin(); err != nil {
        return
    }

    // 3. 结束事务.
    defer func() {
        // 3.1 捕获异常.
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", err)
            log.Panicfc(ctx, "[transaction] fatal: %v.", err)
        }

        // 3.2 结束事务.
        if err != nil {
            e := sess.Rollback()
            if log.Config.DebugOn() {
                if e != nil {
                    log.Debugfc(ctx, "[transaction] rollback: %v.", e)
                } else {
                    log.Debugfc(ctx, "[transaction] rollback.")
                }
            }
        } else {
            e := sess.Commit()
            if log.Config.DebugOn() {
                if e != nil {
                    log.Debugfc(ctx, "[transaction] commit: %v.", e)
                } else {
                    log.Debugfc(ctx, "[transaction] commit.")
                }
            }
        }
    }()

    // 4. 事务过程.
    //    遍历全部执行器, 任一执行器返回error时退出的行.
    for _, handler := range handlers {
        if err = handler(ctx, sess); err != nil {
            break
        }
    }
    return
}
