// author: wsfuyibing <websearch@163.com>
// date: 2022-10-23

package db

import (
	"context"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
	mappers "xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
)

var Connector Connection

type (
	// Connection
	// 连接管理接口.
	Connection interface {
		// GetEngine
		// 读取引擎.
		GetEngine(names ...string) *xorm.EngineGroup

		// GetMaster
		// 读取主库连接.
		GetMaster(names ...string) *xorm.Session

		// GetMasterWithContext
		// 读取主库连接.
		GetMasterWithContext(ctx context.Context, names ...string) *xorm.Session

		// GetSlave
		// 读取从库连接.
		GetSlave(names ...string) *xorm.Session

		// GetSlaveWithContext
		// 读取从库连接.
		GetSlaveWithContext(ctx context.Context, names ...string) *xorm.Session
	}

	// 连接管理.
	connection struct {
		empty   *xorm.EngineGroup
		engines map[string]*xorm.EngineGroup
		mu      sync.RWMutex
	}
)

// GetEngine
// 读取连接.
func (o *connection) GetEngine(names ...string) *xorm.EngineGroup {
	var (
		database *Database
		eg       *xorm.EngineGroup
		exists   bool
		name     = DefaultEngineName
	)

	// 2. 连接名称.
	//    在配置文件 config/db.yaml 中定义.
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	// 1. 读/写锁.
	o.mu.Lock()
	defer o.mu.Unlock()

	// 3. 连接复用.
	if eg = func() *xorm.EngineGroup {
		if v, ok := o.engines[name]; ok {
			return v
		}
		return nil
	}(); eg != nil {
		return eg
	}

	// 2. 读取配置.
	if database, exists = Config.Databases[name]; !exists {
		return o.empty
	}

	// 3. 创建连接.
	eg, _ = xorm.NewEngineGroup(database.Driver, database.Dsn)
	eg.AddHook(&hook{})
	eg.EnableSessionID(*Config.EnableSession)
	eg.SetLogger(Logger)
	eg.ShowSQL(*database.ShowSQL)

	// 3.1 映射关系.
	switch strings.ToLower(database.Mapper) {
	case "gonic":
		eg.SetMapper(mappers.GonicMapper{})
	case "prefix":
		eg.SetMapper(mappers.PrefixMapper{})
	case "same":
		eg.SetMapper(mappers.SameMapper{})
	case "snake":
		eg.SetMapper(mappers.SnakeMapper{})
	case "suffix":
		eg.SetMapper(mappers.SuffixMapper{})
	}

	// 3.2 连接参数.
	eg.SetMaxIdleConns(database.MaxIdle)
	eg.SetMaxOpenConns(database.MaxOpen)
	eg.SetConnMaxLifetime(time.Duration(database.MaxLifetime) * time.Second)

	// 3.3 更新内存.
	o.engines[name] = eg
	return eg
}

// GetMaster
// 读取主库连接.
func (o *connection) GetMaster(names ...string) *xorm.Session {
	return o.GetEngine(names...).
		Master().
		NewSession()
}

// GetMasterWithContext
// 读取主库连接.
func (o *connection) GetMasterWithContext(ctx context.Context, names ...string) *xorm.Session {
	if ctx == nil {
		return o.GetMaster(names...)
	}
	s := o.GetEngine(names...).Master().NewSession()
	s.Context(ctx)
	return s
}

// GetSlave
// 读取从库连接.
func (o *connection) GetSlave(names ...string) *xorm.Session {
	return o.GetEngine(names...).
		Slave().
		NewSession()
}

// GetSlaveWithContext
// 读取从库连接.
func (o *connection) GetSlaveWithContext(ctx context.Context, names ...string) *xorm.Session {
	if ctx == nil {
		return o.GetSlave(names...)
	}
	s := o.GetEngine(names...).Slave().NewSession()
	s.Context(ctx)
	return s
}

// 构造实例.
func (o *connection) init() *connection {
	o.empty, _ = xorm.NewEngineGroup(DefaultDriver, []string{DefaultEmptyDsn})
	o.engines = make(map[string]*xorm.EngineGroup)
	o.mu = sync.RWMutex{}
	return o
}
