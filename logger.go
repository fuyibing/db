// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"fmt"
	"github.com/fuyibing/log/v3"
	xo "xorm.io/xorm/log"
)

// 日志操作.
type logger struct {
	xo.DiscardLogger

	key, data, user string
	undefined       bool
}

// AfterSQL
// SQL执行完成后触发.
func (o *logger) AfterSQL(c xo.LogContext) {
	var (
		prefix = fmt.Sprintf("[SQL=%d]", c.ExecuteTime.Milliseconds())
	)

	// 1. ORM参数.
	if s := c.Ctx.Value(xo.SessionIDKey); s != nil {
		prefix += fmt.Sprintf("[XORM=%s|%s|%s|%s]", o.key, o.user, o.data, s)
	} else {
		prefix += fmt.Sprintf("[XORM=%s|%s|%s]", o.key, o.user, o.data)
	}

	// 2. 查询语句.
	if c.Args != nil && len(c.Args) > 0 {
		log.Client.Infofc(c.Ctx, "%s %s, Args: %v", prefix, c.SQL, c.Args)
	} else {
		log.Client.Infofc(c.Ctx, "%s %s", prefix, c.SQL)
	}

	// 2. 记录Err原因.
	if o.undefined {
		log.Client.Errorfc(c.Ctx, "field '%s' not defined in config file: %v", o.key, c.Err)
	} else if c.Err != nil {
		log.Client.Errorfc(c.Ctx, fmt.Sprintf("%v", c.Err))
	}
}

func (o *logger) BeforeSQL(_ xo.LogContext) {}
func (o *logger) Level() xo.LogLevel        { return xo.LOG_INFO }
func (o *logger) SetLevel(_ xo.LogLevel)    {}
func (o *logger) IsShowSQL() bool           { return true }

func (o *logger) init() *logger { return o }
