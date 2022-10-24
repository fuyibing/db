// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package tests

import (
	"github.com/fuyibing/db/v3"
	"github.com/fuyibing/log/v3"
	"github.com/fuyibing/log/v3/trace"
	"testing"
	"xorm.io/xorm"
)

type (
	Example struct {
		Id int `xorm:"id pk autoincr"`
	}

	ExampleService struct {
		db.Service
	}
)

func (o *Example) TableName() string {
	return "task"
}

func NewExampleService(x ...*xorm.Session) *ExampleService {
	o := &ExampleService{}
	o.Use(x...)
	return o
}

func (o *ExampleService) GetById(id int) (bean *Example, err error) {
	var exists bool
	bean = &Example{}
	if exists, err = o.Slave().Where("id = ?", id).Get(bean); err != nil || !exists {
		bean = nil
	}
	return
}

func TestService(t *testing.T) {
	ctx := trace.New()
	log.Infofc(ctx, "create context")

	sess := db.Connector.GetSlaveWithContext(ctx)
	service := NewExampleService(sess)

	g, err := service.GetById(1)
	if err != nil {
		return
	}
	log.Infofc(ctx, "service completed: %+v", g)
}