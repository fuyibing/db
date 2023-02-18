// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package services

import (
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/db/v8/examples/models"
	"xorm.io/xorm"
)

type ExampleService struct {
	db.Service
}

func NewExampleService(ss ...*xorm.Session) *ExampleService {
	o := &ExampleService{}
	o.Use(ss...)
	o.UseConnection("db") // optional, default value is: db
	return o
}

func (o *ExampleService) GetById(id int) (*models.Example, error) {
	var (
		bean   = &models.Example{}
		err    error
		exists bool
	)
	if exists, err = o.Slave().Where("id = ?", id).Get(bean); err != nil && !exists {
		return nil, err
	}
	return bean, nil
}
