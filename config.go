// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"io/ioutil"
	"time"

	"github.com/fuyibing/log/v2/plugins"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/fuyibing/log/v2"
)

// DB配置.
type configuration struct {
	Driver      string   `yaml:"driver"`
	Dsn         []string `yaml:"dsn"`
	MaxIdle     int      `yaml:"max-idle"`
	MaxOpen     int      `yaml:"max-open"`
	MaxLifetime int      `yaml:"max-lifetime"`
	Mapper      string   `yaml:"mapper"`
	engines     *xorm.EngineGroup
	slaveEnable bool
}

// 读取YAML文件.
func (o *configuration) LoadYaml(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, o); err != nil {
		return err
	}

	log.Debugf("配置源: %s.", path)
	o.parse()
	return nil
}

// 初始化配置.
func (o *configuration) initialize() {
	for _, path := range []string{"./tmp/db.yaml", "./config/db.yaml", "../config/db.yaml"} {
		if o.LoadYaml(path) == nil {
			break
		}
	}
}

// 应用配置.
func (o *configuration) parse() {
	var err error
	if o.engines, err = xorm.NewEngineGroup(o.Driver, o.Dsn); err != nil {
		panic(err)
	}

	log.Debugf("数据源: %d个, 最大连结数: %d个, 表映射: %s.", len(o.Dsn), o.MaxOpen, o.Mapper)

	o.engines.SetConnMaxLifetime(time.Duration(o.MaxLifetime) * time.Second)
	o.engines.SetMaxIdleConns(o.MaxIdle)
	o.engines.SetMaxOpenConns(o.MaxOpen)
	o.engines.SetLogger(plugins.NewXOrm())
	o.slaveEnable = len(o.Dsn) > 1

	if o.Mapper == "same" {
		o.engines.SetColumnMapper(names.SameMapper{})
	} else if o.Mapper == "snake" {
		o.engines.SetColumnMapper(names.SnakeMapper{})
	}
}
