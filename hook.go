// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"context"
	"xorm.io/xorm/contexts"
)

type (
	Hook interface {
	}

	hook struct {
	}
)

func (o *hook) AfterProcess(c *contexts.ContextHook) error {
	return nil
}

func (o *hook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return c.Ctx, nil
}
