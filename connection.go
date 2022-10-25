// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"context"
	"sync"
	"time"
	"xorm.io/xorm"
)

var (
	Connector *Connection
)

type (
	// Connection
	// 连接操作.
	Connection struct {
		engines map[string]*xorm.EngineGroup
		mu      sync.RWMutex
	}
)

// GetEngineGroup
// 读取ORM引擎.
func (o *Connection) GetEngineGroup(keys ...string) (engine *xorm.EngineGroup) {
	var (
		key = o.key(keys...)
		ok  bool
	)

	o.mu.Lock()
	defer o.mu.Unlock()

	if engine, ok = o.engines[key]; !ok {
		engine = o.build(key)
		o.engines[key] = engine
	}

	return
}

// GetMaster
// 读取主库连接.
func (o *Connection) GetMaster(keys ...string) *xorm.Session {
	engine := o.GetEngineGroup(keys...)
	return engine.Master().NewSession()
}

// GetMasterWithContext
// 读取主库连接.
func (o *Connection) GetMasterWithContext(ctx context.Context, keys ...string) *xorm.Session {
	sess := o.GetMaster(keys...)
	sess.Context(ctx)
	return sess
}

// GetSlave
// 读取从库连接.
func (o *Connection) GetSlave(keys ...string) *xorm.Session {
	engine := o.GetEngineGroup(keys...)
	return engine.Slave().NewSession()
}

// GetSlaveWithContext
// 读取从库连接.
func (o *Connection) GetSlaveWithContext(ctx context.Context, keys ...string) *xorm.Session {
	sess := o.GetSlave(keys...)
	sess.Context(ctx)
	return sess
}

// 构建引擎.
func (o *Connection) build(key string) *xorm.EngineGroup {
	var (
		cfg = Config.GetDatabase(key)
		eng *xorm.EngineGroup
	)

	if cfg == nil {
		cfg = Config.GetDefault()
	}

	eng, _ = xorm.NewEngineGroup(cfg.Driver, cfg.Dsn)
	eng.EnableSessionID(*cfg.EnableSession)
	eng.SetMapper(cfg.GetMapper())
	eng.ShowSQL(*cfg.ShowSQL)
	eng.SetMaxIdleConns(cfg.MaxIdle)
	eng.SetMaxOpenConns(cfg.MaxOpen)
	eng.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	eng.SetLogger((&logger{
		data:      cfg.GetDataName(),
		key:       key,
		user:      cfg.GetUsername(),
		undefined: cfg.undefined,
	}).init())

	return eng
}

// 构造实例.
func (o *Connection) init() *Connection {
	o.engines = make(map[string]*xorm.EngineGroup)
	o.mu = sync.RWMutex{}
	return o
}

// 配置名称.
func (o *Connection) key(keys ...string) string {
	if len(keys) > 0 && keys[0] != "" {
		return keys[0]
	}
	return defaultEngine
}
