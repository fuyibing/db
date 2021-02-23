// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"xorm.io/xorm"
)

// Service匿名结构.
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

// 获取主库连结.
func (o *Service) Master() *xorm.Session {
	if o.sess == nil {
		return Master()
	}
	return o.sess
}

// 获取从库连结.
func (o *Service) Slave() *xorm.Session {
	if o.sess == nil {
		return Slave()
	}
	return o.sess
}

// 使用指定连结.
func (o *Service) Use(s ...*xorm.Session) {
	if s != nil && len(s) > 0 {
		o.sess = s[0]
	}
}
