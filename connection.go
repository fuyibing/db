// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"context"
	"sync"
	"time"
	"xorm.io/xorm"
)

var Connector *Connection

type (
	Connection struct {
		engines map[string]*xorm.EngineGroup
		mu      sync.RWMutex
	}
)

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

func (o *Connection) GetMaster(keys ...string) *xorm.Session {
	engine := o.GetEngineGroup(keys...)
	return engine.Master().NewSession()
}

func (o *Connection) GetMasterWithContext(ctx context.Context, keys ...string) *xorm.Session {
	sess := o.GetMaster(keys...)
	sess.Context(ctx)
	return sess
}

func (o *Connection) GetSlave(keys ...string) *xorm.Session {
	engine := o.GetEngineGroup(keys...)
	return engine.Slave().NewSession()
}

func (o *Connection) GetSlaveWithContext(ctx context.Context, keys ...string) *xorm.Session {
	sess := o.GetSlave(keys...)
	sess.Context(ctx)
	return sess
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

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
