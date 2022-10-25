// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"xorm.io/xorm"
)

// Service
// 服务操作.
type Service struct {
	connection string        // 连接名称(默认: db)
	session    *xorm.Session // 自定义连接
}

// Master
// 读取主库连接.
func (o *Service) Master() *xorm.Session {
	if o.session != nil {
		return o.session
	}
	return Connector.GetMaster(o.connection)
}

// Slave
// 读取从库连接.
func (o *Service) Slave() *xorm.Session {
	if o.session != nil {
		return o.session
	}
	return Connector.GetSlave(o.connection)
}

// Use
// 绑定连接.
func (o *Service) Use(sessions ...*xorm.Session) {
	if len(sessions) > 0 {
		o.session = sessions[0]
	}
}

// UseConnection
// 绑定连接名称.
func (o *Service) UseConnection(connection string) {
	o.connection = connection
}
