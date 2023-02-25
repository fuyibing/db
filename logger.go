// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	xo "xorm.io/xorm/log"

	"github.com/fuyibing/log/v5"
)

type logger struct {
	xo.DiscardLogger

	key, data, user string
	undefined       bool
}

func (o *logger) AfterSQL(c xo.LogContext) {
	var (
		span, exists = log.Manager.GetSpan(c.Ctx)
	)

	// Append query statement.
	if exists {
		span.InfoWith(map[string]interface{}{
			"sql-args":     c.Args,
			"sql-duration": c.ExecuteTime.Milliseconds(),
			"sql-session":  c.Ctx.Value(xo.SessionIDKey),
			"sql-config":   o.key,
			"sql-username": o.user,
			"sql-database": o.data,
		}, c.SQL)
	} else {
		log.Info("[SQL] %s, Args: %v", c.SQL, c.Args)
	}

	if o.undefined {
		if exists {
			span.Error("field '%s' not defined in config file: %v", o.key, c.Err)
		} else {
			log.Error("field '%s' not defined in config file: %v", o.key, c.Err)
		}
	} else if c.Err != nil {
		if exists {
			span.Error("%v", c.Err)
		} else {
			log.Error("%v", c.Err)
		}
	}
}

func (o *logger) BeforeSQL(_ xo.LogContext) {}
func (o *logger) Level() xo.LogLevel        { return xo.LOG_INFO }
func (o *logger) SetLevel(_ xo.LogLevel)    {}
func (o *logger) IsShowSQL() bool           { return true }

func (o *logger) init() *logger { return o }
