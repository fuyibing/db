// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	"regexp"
	"strings"
	"xorm.io/xorm/names"
)

var (
	defaultEnableSession = true
	defaultRegexDsn      = regexp.MustCompile(`^([_a-zA-Z0-9]+):[.]*@tcp\([^)]+\)/([^?]+)`)
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

type (
	// Database
	// 数据源定义.
	Database struct {
		Driver        string   `yaml:"driver"`
		Dsn           []string `yaml:"dsn"`
		EnableSession *bool    `yaml:"enable-session"`
		MaxIdle       int      `yaml:"max-idle"`
		MaxLifetime   int      `yaml:"max-lifetime"`
		MaxOpen       int      `yaml:"max-open"`
		Mapper        string   `yaml:"mapper"`
		ShowSQL       *bool    `yaml:"show-sql"`

		data, host, user string
		internal         bool
		mapper           names.Mapper
	}
)

// GetDataName
// 获取数据库名.
func (o *Database) GetDataName() string { return o.data }

// GetHost
// 获取数据源地址.
func (o *Database) GetHost() string { return o.host }

// GetMapper
// 获取映射名.
func (o *Database) GetMapper() names.Mapper { return o.mapper }

// GetUsername
// 获取用户名.
func (o *Database) GetUsername() string { return o.user }

// Internal
// 是否内部源.
func (o *Database) Internal() bool { return o.internal }

// init
// 构造数据源.
func (o *Database) init() *Database {
	if o.Driver == "" {
		o.Driver = defaultDriver
	}
	if o.Mapper == "" {
		o.Mapper = defaultMapperName
	}
	if o.Dsn == nil {
		o.Dsn = make([]string, 0)
	}

	if o.MaxIdle == 0 {
		o.MaxIdle = defaultMaxIdle
	}
	if o.MaxOpen == 0 {
		o.MaxOpen = defaultMaxOpen
	}
	if o.MaxLifetime == 0 {
		o.MaxLifetime = defaultMaxLifetime
	}

	if o.ShowSQL == nil {
		o.ShowSQL = &defaultShowSQL
	}
	if o.EnableSession == nil {
		o.EnableSession = &defaultEnableSession
	}

	// 解析DNS.
	// 从数据源中获取用户, 地址, 库名.
	for _, s := range o.Dsn {
		if m := defaultRegexDsn.FindStringSubmatch(s); len(m) == 4 {
			o.user = m[1]
			o.host = m[2]
			o.data = m[3]
			break
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
