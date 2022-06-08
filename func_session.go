// author: wsfuyibing <websearch@163.com>
// date: 2022-06-07

package db

import "xorm.io/xorm"

// Master
// 主库连接.
//
// 获取默认DB实例的主库连接.
func Master() *xorm.Session {
    sess := Manager.GetEngine().Master().NewSession()
    return sess
}

// MasterContext
// 主库连接.
//
// 获取默认DB实例的主库连接.
func MasterContext(ctx interface{}) *xorm.Session {
    sess := Master()
    sess.Context(Manager.Context(ctx))
    return sess
}

// Slave
// 从库连接.
//
// 获取默认DB实例的从库连接.
func Slave() *xorm.Session {
    sess := Manager.GetEngine().Slave().NewSession()
    return sess
}

// SlaveContext
// 从库连接.
//
// 获取默认DB实例的从库连接.
func SlaveContext(ctx interface{}) *xorm.Session {
    sess := Slave()
    sess.Context(Manager.Context(ctx))
    return sess
}
