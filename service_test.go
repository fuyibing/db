// author: wsfuyibing <websearch@163.com>
// date: 2022-06-07

package db

import (
    "testing"

    "xorm.io/xorm"
)

// Example
// 数据表.
type Example struct {
    TimeZoneId     int64  `xorm:"time_zone_id pk autoincr"`
    UseLeapSeconds string `xorm:"use_leap_seconds"`
}

// TableName
// 返回表名.
func (o *Example) TableName() string {
    return "time_zone"
}

// ExampleService
// 表数据操作.
type ExampleService struct {
    Service
}

// NewExampleService
// 创建表操作实例.
func NewExampleService(s ...*xorm.Session) *ExampleService {
    o := &ExampleService{}
    o.Use(s...)
    return o
}

// GetById
// 按ID读取.
func (o *ExampleService) GetById(id int64) (bean *Example, err error) {
    bean = &Example{}
    _, err = o.Slave().Where("time_zone_id = ?", id).Get(bean)
    return
}

func TestService(t *testing.T) {
    service := NewExampleService()
    bean, err := service.GetById(1000)
    if err != nil {
        t.Errorf("%v", err)
        return
    }
    t.Logf("bean: %+v", bean)
}
