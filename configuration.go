// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var (
	// Config
	// 应用于IRIS的配置.
	Config *Configuration
)

type (
	// Configuration
	// 配置参数.
	Configuration struct {
		// 数据源列表.
		Databases map[string]*Database `yaml:"databases"`

		mu                sync.RWMutex
		undefinedDatabase *Database
	}
)

// GetDatabase
// 按Key名称读取数据源.
func (o *Configuration) GetDatabase(key string) *Database {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if database, ok := o.Databases[key]; ok {
		return database
	}
	return nil
}

// GetUndefined
// 读取数据源.
func (o *Configuration) GetUndefined() *Database {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.undefinedDatabase
}

// SetDatabase
// 设置数据源.
func (o *Configuration) SetDatabase(key string, database *Database) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if database == nil {
		delete(o.Databases, key)
	} else {
		o.Databases[key] = database.init()
	}
}

func (o *Configuration) defaults() {
	if o.Databases == nil {
		o.Databases = make(map[string]*Database)
	}

	for _, v := range o.Databases {
		v.init()
	}
}

func (o *Configuration) init() *Configuration {
	o.mu = sync.RWMutex{}
	o.scan()
	o.defaults()

	o.undefinedDatabase = (&Database{
		Dsn:      []string{defaultEngineDsn},
		internal: true,
	}).init()

	return o
}

func (o *Configuration) scan() {
	for _, f := range []string{"./tmp/db.yaml", "./config/db.yaml", "../tmp/db.yaml", "../config/db.yaml"} {
		if body, err := os.ReadFile(f); err == nil {
			if yaml.Unmarshal(body, o) == nil {
				return
			}
		}
	}
}
