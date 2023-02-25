// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"fmt"
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

	prefix := fmt.Sprintf("[SQL=%d]", c.ExecuteTime.Milliseconds())

	// Append xorm session
	// on log.
	if s := c.Ctx.Value(xo.SessionIDKey); s != nil {
		prefix += fmt.Sprintf("[XORM=%s|%s|%s|%s]", o.key, o.user, o.data, s)
	} else {
		prefix += fmt.Sprintf("[XORM=%s|%s|%s]", o.key, o.user, o.data)
	}

	// Append query statement.
	if c.Args != nil && len(c.Args) > 0 {
		if exists {
			span.Info("%s %s, Args: %v", prefix, c.SQL, c.Args)
		} else {
			log.Info("%s %s, Args: %v", prefix, c.SQL, c.Args)
		}
	} else {
		if exists {
			span.Info("%s %s", prefix, c.SQL)
		} else {
			log.Info("%s %s", prefix, c.SQL)
		}
	}

	if o.undefined {
		if exists {
			span.Error("field '%s' not defined in config file: %v", o.key, c.Err)
		} else {
			log.Error("field '%s' not defined in config file: %v", o.key, c.Err)
		}
	} else if c.Err != nil {
		if exists {
			span.Error(fmt.Sprintf("%v", c.Err))
		} else {
			log.Error(fmt.Sprintf("%v", c.Err))
		}
	}
}

func (o *logger) BeforeSQL(_ xo.LogContext) {}
func (o *logger) Level() xo.LogLevel        { return xo.LOG_INFO }
func (o *logger) SetLevel(_ xo.LogLevel)    {}
func (o *logger) IsShowSQL() bool           { return true }

func (o *logger) init() *logger { return o }
