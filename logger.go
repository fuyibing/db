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
		spa, exists = log.Span(c.Ctx)
	)

	// Append query statement.
	if exists {
		spa.Logger().
			Add("sql-args", c.Args).
			Add("sql-collector", o.key).
			Add("sql-database", o.data).
			Add("sql-duration-ms", c.ExecuteTime.Milliseconds()).
			Add("sql-session", c.Ctx.Value(xo.SessionIDKey)).
			Add("sql-username", o.user).
			Info(c.SQL)
	} else {
		log.Info("[SQL] %s, Args: %v", c.SQL, c.Args)
	}

	if o.undefined {
		if exists {
			spa.Logger().Error("field '%s' not defined in config file: %v", o.key, c.Err)
		} else {
			log.Error("field '%s' not defined in config file: %v", o.key, c.Err)
		}
	} else if c.Err != nil {
		if exists {
			spa.Logger().Error("%v", c.Err)
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
