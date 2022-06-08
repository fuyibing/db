// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

// Package db used for database operations.
//
//
package db

import "sync"

func init() {
    new(sync.Once).Do(func() {
        // 1. 单例实例.
        //    - 配置实例.
        //    - 管理器实例.
        Config = (&Configuration{}).Init()
        Manager = (&Management{}).Init()
    })
}
