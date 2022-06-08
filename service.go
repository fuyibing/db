// author: wsfuyibing <websearch@163.com>
// date: 2022-06-07

package db

import "xorm.io/xorm"

// Service
// 服务结构体.
type Service struct {
    _key  string
    _sess *xorm.Session
}

// Master
// 获取主库连接.
func (o *Service) Master() *xorm.Session {
    // 1. 指定连接.
    if o._sess != nil {
        return o._sess
    }

    // 2. 创建连接.
    return Manager.GetEngine(o._key).Master().NewSession()
}

// Slave
// 获取从库连接.
func (o *Service) Slave() *xorm.Session {
    // 1. 指定连接.
    if o._sess != nil {
        return o._sess
    }

    // 2. 创建连接.
    return Manager.GetEngine(o._key).Slave().NewSession()
}

// Use
// 绑定连接.
func (o *Service) Use(s ...*xorm.Session) {
    if len(s) > 0 {
        o._sess = s[0]
    }
}

// UseKey
// 绑定连接键名.
//
// 键名用于选定连接, 以支持在一个项目中可以连接多个实例.
func (o *Service) UseKey(key string) {
    if key != "" {
        o._key = key
    }
}
