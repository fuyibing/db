// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package db

import (
	"sync"

	"github.com/fuyibing/log"
)

var (
	Config *configuration
)

func init() {
	new(sync.Once).Do(func() {
		log.Info("initialize golang framework service.")
		Config = new(configuration)
		Config.initialize()
	})
}
