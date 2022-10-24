// author: wsfuyibing <websearch@163.com>
// date: 2022-10-24

package tests

import (
	"github.com/fuyibing/db/v3"
	"testing"
)

func TestConfig(t *testing.T) {
	c := db.Config
	t.Logf("config: %+v", c)
}
