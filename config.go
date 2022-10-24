// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"github.com/fuyibing/log/v3"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	// Config
	// 配置.
	Config *config

	// 启用链路.
	//
	// 当启用时, 在SQL执行的日志中, 打印SessionID, 同一个事务
	// 保持相同的SessionID.
	defaultEnableSession = true

	// 映射关系.
	defaultMapper = "snake"

	// 输出SQL.
	defaultShowSQL = true
)

const (
	DefaultDriver      = "mysql"
	DefaultEmptyDsn    = "username:password@tcp(127.0.0.1:3306)/mysql?charset=utf8"
	DefaultEngineName  = "db"
	DefaultMaxIdle     = 2
	DefaultMaxLifetime = 60
	DefaultMaxOpen     = 30
)

// 基础配置.
type config struct {
	// 连接参数.
	Databases map[string]*Database `yaml:"databases"`

	// 链路状态.
	EnableSession *bool `yaml:"enable-session"`
}

// 赋默认值.
func (o *config) defaults() *config {
	// 1. 连结配置.
	//    支持多连接.
	if o.Databases == nil {
		o.Databases = make(map[string]*Database)
	}

	// 1.1 默认参数.
	for k, v := range o.Databases {
		// 驱动名.
		if v.Driver == "" {
			v.Driver = DefaultDriver
		}
		// 最大空闲.
		if v.MaxIdle == 0 {
			v.MaxIdle = DefaultMaxIdle
		}
		// 最大活跃期.
		if v.MaxLifetime == 0 {
			v.MaxLifetime = DefaultMaxLifetime
		}
		// 最大连接数.
		if v.MaxOpen == 0 {
			v.MaxOpen = DefaultMaxOpen
		}
		// 映射配置.
		if v.Mapper == "" {
			v.Mapper = defaultMapper
		}

		if v.ShowSQL == nil {
			v.ShowSQL = &defaultShowSQL
		}
		log.Debugf("database found: driver=%s, name=%s, dsn-item=%d", v.Driver, k, len(v.Dsn))
	}

	// 2. Sess 链路.
	if o.EnableSession == nil {
		o.EnableSession = &defaultEnableSession
	}

	return o
}

// 构造实例.
func (o *config) init() *config {
	log.Debug("database init")
	return o.scan().defaults()
}

// 扫描配置.
func (o *config) scan() *config {
	for _, f := range []string{"./tmp/db.yaml", "./config/db.yaml", "../tmp/db.yaml", "../config/db.yaml"} {
		if body, err := os.ReadFile(f); err == nil {
			if yaml.Unmarshal(body, o) == nil {
				log.Debugf("database load: %v", f)
				break
			}
		}
	}
	return o
}
