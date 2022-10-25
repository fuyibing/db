// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

// Package db
// 基于 XORM 的数据库操作工具.
package db

import (
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		Config = (&Configuration{}).init()
		Connector = (&Connection{}).init()
	})
}
