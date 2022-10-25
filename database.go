// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"regexp"
	"strings"
	"xorm.io/xorm/names"
)

var (
	defaultEnableSession = true
	defaultShowSQL       = true
)

const (
	defaultDriver = "mysql"

	defaultEngine    = "db"
	defaultEngineDsn = "user:pass@tcp(host:3306)/mysql?charset=utf8"

	defaultMaxIdle     = 2
	defaultMaxLifetime = 60
	defaultMaxOpen     = 30

	defaultMapperName = "snake"
)

// Database
// 数据库配置.
type Database struct {
	Driver        string   `yaml:"driver"`
	Dsn           []string `yaml:"dsn"`
	MaxIdle       int      `yaml:"max-idle"`
	MaxLifetime   int      `yaml:"max-lifetime"`
	MaxOpen       int      `yaml:"max-open"`
	Mapper        string   `yaml:"mapper"`
	EnableSession *bool    `yaml:"enable-session"`
	ShowSQL       *bool    `yaml:"show-sql"`

	data      string
	mapper    names.Mapper
	undefined bool
	user      string
}

// GetMapper
// 返回映射关系.
func (o *Database) GetMapper() names.Mapper {
	return o.mapper
}

// GetDataName
// 返回数据库库名.
func (o *Database) GetDataName() string {
	return o.data
}

// GetUsername
// 返回数据库用户名.
func (o *Database) GetUsername() string {
	return o.user
}

// Undefined
// 是否未定义.
func (o *Database) Undefined() bool {
	return o.undefined
}

// 初始化.
func (o *Database) init() *Database {
	// 基础项.

	if o.Driver == "" {
		o.Driver = defaultDriver
	}
	if o.Mapper == "" {
		o.Mapper = defaultMapperName
	}
	if o.Dsn == nil {
		o.Dsn = make([]string, 0)
	}

	// 连接池

	if o.MaxIdle == 0 {
		o.MaxIdle = defaultMaxIdle
	}
	if o.MaxOpen == 0 {
		o.MaxOpen = defaultMaxOpen
	}
	if o.MaxLifetime == 0 {
		o.MaxLifetime = defaultMaxLifetime
	}

	// 附加项.

	if o.ShowSQL == nil {
		o.ShowSQL = &defaultShowSQL
	}
	if o.EnableSession == nil {
		o.EnableSession = &defaultEnableSession
	}

	// 数据库.

	for _, s := range o.Dsn {
		if m := regexp.MustCompile(`^([_a-zA-Z0-9\-]+):([^/]+)/([_a-zA-Z0-9\-]+)`).FindStringSubmatch(s); len(m) > 0 {
			o.data = m[3]
			o.user = m[1]
		}
	}

	// 映射关系.

	switch strings.ToLower(o.Mapper) {
	case "gonic":
		o.mapper = names.GonicMapper{}
	case "prefix":
		o.mapper = names.PrefixMapper{}
	case "same":
		o.mapper = names.SameMapper{}
	case "snake":
		o.mapper = names.SnakeMapper{}
	case "suffix":
		o.mapper = names.SuffixMapper{}
	}

	return o
}
