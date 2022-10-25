// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package db

import (
	"testing"
)

func ExampleConfiguration_SetDatabase() {
	Config.SetDatabase("my", &Database{
		Dsn: []string{
			"user:pass@tcp(192.168.0.100:6379)/schema?charset=utf8",
			"user:pass@tcp(192.168.0.101:6379)/schema?charset=utf8",
			"user:pass@tcp(192.168.0.102:6379)/schema?charset=utf8",
			"user:pass@tcp(192.168.0.103:6379)/schema?charset=utf8",
		},
	})
}

func TestConfiguration_SetDatabase(t *testing.T) {
	ExampleConfiguration_SetDatabase()
}
