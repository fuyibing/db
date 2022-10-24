// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package db

import (
	"xorm.io/xorm"
)

// Service
// ORM数据服务.
//
//   type ExampleService struct {
//       Service
//   }
type Service struct {
	// 连接名称.
	// 系统基于此名称, 加载连接选项.
	_name string

	// 连接实例.
	// 从连接池中获取活跃连接.
	_sess *xorm.Session
}

// Master
// 读取主库连接.
func (o *Service) Master() *xorm.Session {
	if o._sess != nil {
		return o._sess
	}
	return Connector.GetMaster(o._name)
}

// Slave
// 读取从库连接.
func (o *Service) Slave() *xorm.Session {
	if o._sess != nil {
		return o._sess
	}
	return Connector.GetSlave(o._name)
}

// Use
// 使用指定连接.
func (o *Service) Use(s ...*xorm.Session) {
	if len(s) > 0 {
		o._sess = s[0]
	}
}

// UseConnection
// 使用连接名称.
func (o *Service) UseConnection(name string) {
	o._name = name
}
