// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

package db

import (
    "os"
    "sync"

    "gopkg.in/yaml.v3"
)

// Config
// 配置实例/单例.
var Config *Configuration

// Configuration
// 配置结构体.
type Configuration struct {
    Databases    map[string]*Database `yaml:"databases"`
    Mapper       *string              `yaml:"mapper"`
    MaxIdle      *int                 `yaml:"max-idle"`
    MaxLifetime  *int                 `yaml:"max-lifetime"`
    MaxOpen      *int                 `yaml:"max-open"`
    UseSessionId *bool                `yaml:"use-session-id"`

    mu *sync.RWMutex
}

// GetDatabase
// 读取Database实例.
func (o *Configuration) GetDatabase(key string) *Database {
    o.mu.RLock()
    defer o.mu.RUnlock()
    if v, ok := o.Databases[key]; ok {
        return v
    }
    return nil
}

// Init
// 配置初始化.
func (o *Configuration) Init() *Configuration {
    o.mu = new(sync.RWMutex)
    o.Databases = make(map[string]*Database)

    o.scan()
    o.defaults()
    return o
}

// SetDatabase
// 设置Database实例.
func (o *Configuration) SetDatabase(key string, database *Database) {
    o.mu.Lock()
    defer o.mu.Unlock()
    o.Databases[key] = database
}

// 赋默认值.
func (o *Configuration) defaults() {
    // 1. 映射模式.
    if o.Mapper == nil {
        o.Mapper = &defaultMapper
    }

    // 2. 最大空闲.
    if o.MaxIdle == nil {
        o.MaxIdle = &defaultMaxIdle
    }

    // 3. 连接时长.
    //    连接在池中超过指定时长没有被使用后, 关闭连接.
    if o.MaxLifetime == nil {
        o.MaxLifetime = &defaultMaxLifetime
    }

    // 4. 打开文件.
    if o.MaxOpen == nil {
        o.MaxOpen = &defaultMaxOpen
    }

    // 5. 连接标识.
    if o.UseSessionId == nil {
        o.UseSessionId = &defaultUseSessionId
    }
}

// 扫描配置.
// 扫描 db 配置文件(db.yaml).
func (o *Configuration) scan() {
    for _, f := range []string{"tmp/db.yaml", "config/db.yaml", "../tmp/db.yaml", "../config/db.yaml"} {
        // 1. 读取文件.
        buf, err := os.ReadFile(f)

        // 2. 读取出错.
        //    文件不存在或无读取权限.
        if err != nil {
            continue
        }

        // 3. 格式匹配.
        //    任意一次解析正确后, 退出执行.
        if yaml.Unmarshal(buf, o) == nil {
            break
        }
    }
}
