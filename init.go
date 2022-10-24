// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

// Package db used for database operations.
//
//
package db

import "sync"

func init() {
	new(sync.Once).Do(func() {
		Config = (&config{}).init()
		Logger = (&logger{}).init()
		Connector = (&connection{}).init()
	})
}
