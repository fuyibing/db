// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package main

import (
	"github.com/fuyibing/db/v8"
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
	ctx := log.NewContextInfo("trace begin")
	sess := db.Connector.GetMasterWithContext(ctx)

	defer func() {
		log.Infofc(ctx, "trace end")
		log.Client.Stop()
	}()

	bean, err := services.NewExampleService(sess).GetById(1)
	if err != nil {
		return
	}

	log.Infofc(ctx, "found: %d", bean.Id)
}
