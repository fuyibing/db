// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package db

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var Config *Configuration

type (
	Configuration struct {
		Databases map[string]*Database `yaml:"databases"`

		database *Database
		mu       sync.RWMutex
	}
)

func (o *Configuration) GetDatabase(key string) *Database {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if database, ok := o.Databases[key]; ok {
		return database
	}
	return nil
}

func (o *Configuration) GetDefault() *Database { return o.database }

func (o *Configuration) SetDatabase(key string, database *Database) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if database == nil {
		delete(o.Databases, key)
	} else {
		o.Databases[key] = database.init()
	}
}

// /////////////////////////////////////////////////////////////
// Access methods
// /////////////////////////////////////////////////////////////

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

	o.database = (&Database{Dsn: []string{defaultEngineDsn}, undefined: true}).init()
	return o
}

func (o *Configuration) scan() {
	for _, f := range []string{"./tmp/db.yaml", "./config/db.yaml", "../tmp/db.yaml", "../config/db.yaml"} {
		if body, err := os.ReadFile(f); err == nil {
			if yaml.Unmarshal(body, o) == nil {
				break
			}
		}
	}
}
