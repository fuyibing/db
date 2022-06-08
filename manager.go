// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

package db

import (
    "context"
    "sync"
    "time"

    "github.com/fuyibing/log/v2"
    "github.com/fuyibing/log/v2/interfaces"
    "github.com/fuyibing/log/v2/plugins"
    "github.com/kataras/iris/v12"
    "xorm.io/xorm"
    "xorm.io/xorm/names"

    _ "github.com/go-sql-driver/mysql"
)

// Manager
// 管理器实例/单例.
var Manager *Management

// Management
// 管理器结构体.
type Management struct {
    engines map[string]*xorm.EngineGroup
    mu      *sync.RWMutex
}

// Context
// 添加上下文.
func (o *Management) Context(i interface{}) context.Context {
    // 1. 空上下文.
    //    若入参为空时, 基于log包创建默认上下文.
    if i == nil {
        return log.NewContext()
    }

    // 2. 类型转换.
    if ctx, ok := i.(context.Context); ok {
        return ctx
    }

    // 3. 三方框架.
    //    基于iris框架构建的上下文.
    if c, ok := i.(iris.Context); ok {
        if g := c.Values().Get(interfaces.OpenTracingKey); g != nil {
            return context.WithValue(context.Background(), interfaces.OpenTracingKey, g.(interfaces.TraceInterface))
        }
    }

    // 4. 丢弃入参.
    return log.NewContext()
}

// GetEngine
// 读取DB引擎.
func (o *Management) GetEngine(keys ...string) (eg *xorm.EngineGroup) {
    // 1. 准备读取.
    var (
        cfg *Database
        err error
        key = defaultDatabaseKey
        ok  = false
    )

    // 2. 默认名称.
    if len(keys) > 0 && keys[0] != "" {
        key = keys[0]
    }

    // 3. 互拆锁.
    o.mu.Lock()
    defer o.mu.Unlock()

    // 4. 复用引擎.
    if eg, ok = o.engines[key]; ok {
        return
    }

    // 5. 读取配置.
    if cfg = Config.GetDatabase(key); cfg == nil {
        log.Errorf("database config not defined: %s.", key)
        return
    }

    // 6. 创建引擎.
    eg, err = xorm.NewEngineGroup(cfg.Driver, cfg.Dsn)
    if err != nil {
        log.Errorf("create database engine error: name=%s, error=%v.", key, err)
        if eg != nil {
            eg = nil
        }
        return
    }

    // 6.1 日志链路.
    eg.EnableSessionID(*Config.UseSessionId)
    eg.SetLogger(plugins.NewXOrm())

    // 6.2 连接参数.
    eg.SetMaxIdleConns(*Config.MaxIdle)
    eg.SetConnMaxLifetime(time.Duration(*Config.MaxLifetime) * time.Second)
    eg.SetMaxOpenConns(*Config.MaxOpen)

    // 6.3 映射关系.
    switch *Config.Mapper {
    case "snake":
        eg.SetMapper(names.SnakeMapper{})
    case "same":
        eg.SetMapper(names.SameMapper{})
    }

    // 7. 更新内存.
    o.engines[key] = eg
    return
}

// Init
// 管理器初始化.
func (o *Management) Init() *Management {
    o.engines = make(map[string]*xorm.EngineGroup)
    o.mu = new(sync.RWMutex)
    return o
}

// SetEngine
// 设置DB引擎.
func (o *Management) SetEngine(key string, eg *xorm.EngineGroup) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.engines[key] = eg
}
