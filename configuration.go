// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package db

import (
	"github.com/fuyibing/log/v3"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var (
	// Config
	// 配置实例.
	Config *Configuration
)

type (
	// Configuration
	// 配置结构体.
	Configuration struct {
		// 数据源列表.
		// 包初始化时, 从 config/db.yaml 文件中解析.
		Databases map[string]*Database `yaml:"databases"`

		database *Database    // 默认数据源
		mu       sync.RWMutex // 读写锁
	}
)

// GetDatabase
// 读取数据源.
func (o *Configuration) GetDatabase(key string) *Database {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if database, ok := o.Databases[key]; ok {
		return database
	}
	return nil
}

// GetDefault
// 读取数据源.
func (o *Configuration) GetDefault() *Database {
	return o.database
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

// 默认配置.
func (o *Configuration) defaults() {
	if o.Databases == nil {
		o.Databases = make(map[string]*Database)
	}
	for _, v := range o.Databases {
		v.init()
	}
}

// 构造实例.
func (o *Configuration) init() *Configuration {
	log.Debug("database init")

	o.scan()
	o.defaults()

	o.database = (&Database{Dsn: []string{defaultEngineDsn}, undefined: true}).init()
	o.mu = sync.RWMutex{}
	return o
}

// 扫描配置.
func (o *Configuration) scan() {
	for _, f := range []string{"./tmp/db.yaml", "./config/db.yaml", "../tmp/db.yaml", "../config/db.yaml"} {
		if body, err := os.ReadFile(f); err == nil {
			if yaml.Unmarshal(body, o) == nil {
				log.Debugf("database load: %v", f)
				break
			}
		}
	}
}
