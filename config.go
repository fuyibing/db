// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/fuyibing/log"
)

// configuration.
type configuration struct {
	// Config fields.
	Driver      string   `yaml:"driver"`
	Dsn         []string `yaml:"dsn"`
	MaxIdle     int      `yaml:"max-idle"`
	MaxOpen     int      `yaml:"max-open"`
	MaxLifetime int      `yaml:"max-lifetime"`
	Mapper      string   `yaml:"mapper"`
	// Engine group for xorm.
	engines *xorm.EngineGroup
}

// Load configuration from yaml file.
func (o *configuration) LoadYaml(path string) error {
	data, err := ioutil.ReadFile(path)
	// return if read file error.
	if err != nil {
		return err
	}
	// return if parse yaml error.
	if err = yaml.Unmarshal(data, o); err != nil {
		return err
	}
	// parse config.
	log.Infof("load config from %s.", path)
	o.parse()
	return nil
}

// init config with default file.
func (o *configuration) initialize() {
	for _, path := range []string{"./config/service.yaml", "../config/service.yaml"} {
		err := o.LoadYaml(path)
		if err == nil {
			break
		}
	}
}

// parse marshal.
func (o *configuration) parse() {
	// prepare.
	var err error
	if o.engines, err = xorm.NewEngineGroup(o.Driver, o.Dsn); err != nil {
		panic(err)
	}
	// initialize options.
	log.Infof("assign %s driver with %d dsn, max idles is %d, max open files is %d.", o.Driver, len(o.Dsn), o.MaxIdle, o.MaxOpen)
	o.engines.SetConnMaxLifetime(time.Duration(o.MaxLifetime) * time.Second)
	o.engines.SetMaxIdleConns(o.MaxIdle)
	o.engines.SetMaxOpenConns(o.MaxOpen)
	// use specified log adapter.
	o.engines.SetLogger(log.NewXOrmLog())
	// fields mapping.
	if o.Mapper == "same" {
		o.engines.SetColumnMapper(names.SameMapper{})
	} else if o.Mapper == "snake" {
		o.engines.SetColumnMapper(names.SnakeMapper{})
	}
}
