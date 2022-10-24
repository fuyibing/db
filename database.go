// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

// Database
// 数据库配置.
type Database struct {
	Driver string   `yaml:"driver"`
	Dsn    []string `yaml:"dsn"`

	MaxIdle     int `yaml:"max-idle"`
	MaxLifetime int `yaml:"max-lifetime"`
	MaxOpen     int `yaml:"max-open"`

	ShowSQL *bool  `yaml:"show-sql"`
	Mapper  string `yaml:"mapper"`
}
