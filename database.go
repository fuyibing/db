// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

// Database
// 数据库配置.
type Database struct {
	// 驱动名.
	Driver string `yaml:"driver"`

	// 数据源.
	Dsn []string `yaml:"dsn"`

	// 最大空闲.
	MaxIdle int `yaml:"max-idle"`

	// 生命周期.
	MaxLifetime int `yaml:"max-lifetime"`

	// 最大打开文件.
	MaxOpen int `yaml:"max-open"`

	// 显示SQL.
	ShowSQL *bool `yaml:"show-sql"`

	// 映射关系.
	Mapper string `yaml:"mapper"`
}
