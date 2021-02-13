// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"xorm.io/xorm"
)

// Root service.
//
//   type ExampleService struct{
//       xdb.Service
//   }
//
//   func NewExampleService(s ...*xorm.Session) *ExampleService {
//       o := &ExampleService{}
//       o.Use(s...)
//       return o
//   }
//
type Service struct {
	sess *xorm.Session
}

// Return master connection session.
func (o *Service) Master() *xorm.Session {
	if o.sess == nil {
		return Config.engines.Master().NewSession()
	}
	return o.sess
}

// Return slave connection session.
func (o *Service) Slave() *xorm.Session {
	if o.sess == nil {
		return Config.engines.Slave().NewSession()
	}
	return o.sess
}

// Use connection session.
func (o *Service) Use(s ...*xorm.Session) {
	if s != nil && len(s) > 0 {
		o.sess = s[0]
	}
}
