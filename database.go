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

func (o *Database) GetMapper() names.Mapper { return o.mapper }
func (o *Database) GetDataName() string     { return o.data }
func (o *Database) GetUsername() string     { return o.user }
func (o *Database) Undefined() bool         { return o.undefined }

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

	for _, s := range o.Dsn {
		if m := regexp.MustCompile(`^([_a-zA-Z0-9\-]+):([^/]+)/([_a-zA-Z0-9\-]+)`).FindStringSubmatch(s); len(m) > 0 {
			o.data = m[3]
			o.user = m[1]
		}
	}

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
