// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	xo "xorm.io/xorm/log"

	"github.com/fuyibing/log/v5"
)

type (
	// logger
	// 覆盖XORM日志.
	logger struct {
		xo.DiscardLogger

		key, data, host, user string
		internal              bool
	}
)

// AfterSQL
// 后置事件.
//
// 当SQL执行完成后, XORM调用此方法.
func (o *logger) AfterSQL(c xo.LogContext) {
	spa, exists := log.Span(c.Ctx)

	// 追加日志.
	if exists {
		// 链路模式.
		// 基于 `fuyibing/log` 中间件的 OpenTelemetry 规范.
		spa.Logger().
			Add("sql-cfg", o.key).
			Add("sql-sess", c.Ctx.Value(xo.SessionIDKey)).
			Add("sql-duration-ms", c.ExecuteTime.Milliseconds()).
			Add("sql-arg", c.Args).
			Add("sql-data", o.data).
			Add("sql-user", o.user).
			Add("sql-host", o.host).
			Info(c.SQL)
	} else {
		// 普通日志.
		log.Info("[SQL] %s, Args: %v", c.SQL, c.Args)
	}

	// 执行出错.
	if c.Err != nil {
		if exists {
			// 链路模式.
			spa.Logger().Error("%v", c.Err)
		} else {
			// 普通模式.
			log.Error("%v", c.Err)
		}
	}
}

func (o *logger) BeforeSQL(_ xo.LogContext) {}
func (o *logger) Level() xo.LogLevel        { return xo.LOG_INFO }
func (o *logger) SetLevel(_ xo.LogLevel)    {}
func (o *logger) IsShowSQL() bool           { return true }

func (o *logger) init() *logger { return o }
