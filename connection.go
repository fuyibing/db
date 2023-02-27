// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"context"
	"sync"
	"time"
	"xorm.io/xorm"
)

var (
	// Connector
	// 连接管理.
	Connector *Connection
)

type (
	// Connection
	// 连接管理接口.
	Connection struct {
		engines map[string]*xorm.EngineGroup
		mu      sync.RWMutex
	}
)

// GetEngineGroup
// 读取XORM引擎组.
func (o *Connection) GetEngineGroup(keys ...string) (engine *xorm.EngineGroup) {
	var (
		key = o.key(keys...)
		ok  bool
	)

	o.mu.Lock()
	defer o.mu.Unlock()

	// 复用连接组配置, 若不存在则创建.
	if engine, ok = o.engines[key]; !ok {
		engine = o.build(key)
		o.engines[key] = engine
	}

	return
}

// GetMaster
// 从Master获取Session.
func (o *Connection) GetMaster(keys ...string) *xorm.Session {
	return o.GetEngineGroup(keys...).
		Master().
		NewSession()
}

// GetMasterWithContext
// 基于Context从Master获取Session.
func (o *Connection) GetMasterWithContext(ctx context.Context, keys ...string) *xorm.Session {
	return o.GetMaster(keys...).
		Context(ctx)
}

// GetSlave
// 从Slave获取Session.
func (o *Connection) GetSlave(keys ...string) *xorm.Session {
	return o.GetEngineGroup(keys...).
		Slave().
		NewSession()
}

// GetSlaveWithContext
// 基于Context从Slave获取Session.
func (o *Connection) GetSlaveWithContext(ctx context.Context, keys ...string) *xorm.Session {
	return o.GetSlave(keys...).
		Context(ctx)
}

func (o *Connection) build(key string) *xorm.EngineGroup {
	var (
		cfg = Config.GetDatabase(key)
		eng *xorm.EngineGroup
	)

	if cfg == nil {
		cfg = Config.GetUndefined()
	}

	eng, _ = xorm.NewEngineGroup(cfg.Driver, cfg.Dsn)
	eng.EnableSessionID(*cfg.EnableSession)
	eng.SetMapper(cfg.GetMapper())
	eng.ShowSQL(*cfg.ShowSQL)
	eng.SetMaxIdleConns(cfg.MaxIdle)
	eng.SetMaxOpenConns(cfg.MaxOpen)
	eng.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	eng.SetLogger((&logger{
		data:     cfg.GetDataName(),
		host:     cfg.GetHost(),
		key:      key,
		user:     cfg.GetUsername(),
		internal: cfg.internal,
	}).init())

	return eng
}

func (o *Connection) init() *Connection {
	o.engines = make(map[string]*xorm.EngineGroup)
	o.mu = sync.RWMutex{}
	return o
}

func (o *Connection) key(keys ...string) string {
	if len(keys) > 0 && keys[0] != "" {
		return keys[0]
	}
	return defaultEngine
}
