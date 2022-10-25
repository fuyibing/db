// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package db

import (
	"testing"
	"xorm.io/xorm"
)

type (
	example struct {
	}

	exampleService struct {
		Service
	}
)

func newExampleService(s ...*xorm.Session) *exampleService {
	service := &exampleService{}
	service.Use(s...)
	service.UseConnection("db")
	return service
}

func (o *exampleService) GetById(id int) (*example, error) {
	var (
		bean   = &example{}
		err    error
		exists = false
	)
	if exists, err = o.Slave().Where("`id` = ?", id).Get(bean); err != nil || !exists {
		bean = nil
	}
	return bean, err
}

func ExampleService_Master() {
	bean := &example{}
	service := &exampleService{}
	exists, err := service.Master().Get(bean)

	if err != nil {
		println("service error:", err.Error())
		return
	}

	if exists {
		println("service found")
	} else {
		println("service not found")
	}
}

func ExampleService_Slave() {
	bean := &example{}
	service := &exampleService{}
	service.UseConnection("my")
	exists, err := service.Slave().Get(bean)

	if err != nil {
		println("service error:", err.Error())
		return
	}

	if exists {
		println("service found")
	} else {
		println("service not found")
	}
}

func ExampleService_Use() {
	session := Connector.GetMaster()

	service := &exampleService{}
	service.Use(session)
}

func ExampleService_UseConnection() {
	// [config.yaml]
	//
	// databases:
	//   my:
	//     driver: "mysql"
	//     dsn:
	//       - "user:pass@tcp(192.168.1.100:3306)/schema?charset=utf8" # master
	//       - "user:pass@tcp(192.168.1.101:3306)/schema?charset=utf8" # slave1
	//       - "user:pass@tcp(192.168.1.102:3306)/schema?charset=utf8" # slave2
	//       - "user:pass@tcp(192.168.1.103:3306)/schema?charset=utf8" # slave3
	service := &exampleService{}
	service.UseConnection("my")
}

func TestService_Master(t *testing.T) {
	ExampleService_Master()
}

func TestService_Slave(t *testing.T) {
	ExampleService_Slave()
}

func TestService_Use(t *testing.T) {
	ExampleService_Use()
}

func TestService_UseConnection(t *testing.T) {
	ExampleService_UseConnection()
}
