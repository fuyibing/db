// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"sync"

	"github.com/fuyibing/log/v2"
)

var (
	Config *configuration
)

func init() {
	new(sync.Once).Do(func() {
		log.Debug("init db package.")
		Config = new(configuration)
		Config.initialize()
	})
}
