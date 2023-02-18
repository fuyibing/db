// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"xorm.io/xorm"
)

type Service struct {
	connection string
	session    *xorm.Session
}

func (o *Service) Master() *xorm.Session {
	if o.session != nil {
		return o.session
	}
	return Connector.GetMaster(o.connection)
}

func (o *Service) Slave() *xorm.Session {
	if o.session != nil {
		return o.session
	}
	return Connector.GetSlave(o.connection)
}

func (o *Service) Use(sessions ...*xorm.Session) {
	if len(sessions) > 0 {
		o.session = sessions[0]
	}
}

func (o *Service) UseConnection(connection string) {
	o.connection = connection
}
