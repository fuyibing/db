// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

// Package db
// database plugin based on xorm.
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
