// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"xorm.io/xorm"
)

type (
	// Service
	// 服务封装.
	Service struct {
		// 连接名.
		// 若不指定, 则使用默认值.
		connection string

		// 连接.
		session *xorm.Session
	}
)

// Master
// 使用主库连接.
func (o *Service) Master() *xorm.Session {
	if o.session != nil {
		return o.session
	}

	return Connector.GetMaster(o.connection)
}

// Slave
// 使用从库连接.
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
// 绑定连接名.
func (o *Service) UseConnection(connection string) {
	o.connection = connection
}
