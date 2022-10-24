// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"fmt"
	"github.com/fuyibing/log/v3"
	xlog "xorm.io/xorm/log"
)

var (
	// Logger
	// 日志.
	Logger *logger
)

// 日志结构体.
type logger struct {
	xlog.DiscardLogger
}

// AfterSQL
// 后置日志.
//
// 执行 SQL 完成后, 记录其结果.
func (o *logger) AfterSQL(c xlog.LogContext) {
	sid := "[SQL]"

	// 1. 打印链路.
	if *Config.EnableSession {
		if s := c.Ctx.Value(xlog.SessionIDKey); s != nil {
			sid = fmt.Sprintf("[SQL=%v]", s)
		}
	}

	// 2. 构建语句.
	if c.Args != nil && len(c.Args) > 0 {
		log.Client.Infofc(c.Ctx, fmt.Sprintf("%v[d=%.06f] %s, Params: %v",
			sid, c.ExecuteTime.Seconds(), c.SQL, c.Args,
		))
	} else {
		log.Client.Infofc(c.Ctx, fmt.Sprintf("%v[d=%.06f] %s",
			sid, c.ExecuteTime.Seconds(), c.SQL,
		))
	}

	// 2. 记录Err原因.
	if c.Err != nil {
		log.Client.Errorfc(c.Ctx, fmt.Sprintf("%v%v", sid, c.Err))
	}
}

// BeforeSQL
// 前置日志.
func (o *logger) BeforeSQL(_ xlog.LogContext) {}

func (o *logger) Level() xlog.LogLevel     { return xlog.LOG_INFO }
func (o *logger) SetLevel(_ xlog.LogLevel) {}
func (o *logger) IsShowSQL() bool          { return true }

func (o *logger) init() *logger { return o }
