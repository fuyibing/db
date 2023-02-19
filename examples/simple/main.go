// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package main

import (
	"github.com/fuyibing/db/v8/examples/services"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/log/v8/conf"
)

func init() {
	log.Config.Set(
		conf.SetAdapter("term"),
		conf.SetLevel("DEBUG"),
		conf.SetTermColor(true),
	)
}

func main() {
	defer log.Client.Close()

	bean, err := services.NewExampleService().GetById(1)

	if err != nil {
		return
	}

	log.Infof("found: %d", bean.Id)
}
