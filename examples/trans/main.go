// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package main

import (
	"context"
	"github.com/fuyibing/db/v8"
	"github.com/fuyibing/db/v8/examples/services"
	"github.com/fuyibing/log/v8"
	"github.com/fuyibing/log/v8/base"
	"github.com/fuyibing/log/v8/conf"
	"xorm.io/xorm"
)

type formatter struct{}

func (o *formatter) Body(_ *base.Line) []byte { return nil }
func (o *formatter) String(p *base.Line) string {
	return p.Text
}

func init() {
	log.Config.Set(
		conf.SetAdapter("term"),
		conf.SetLevel("DEBUG"),
		conf.SetTermColor(true),
	)

	log.Client.GetAdapterRegistry().SetFormatter(&formatter{})
}

func main() {
	ctx := log.NewContextInfo("--------------trace begin")
	sess := db.Connector.GetMasterWithContext(ctx)

	defer func() {
		log.Infofc(ctx, "-------------- trace end")
		log.Client.Close()
	}()

	_ = db.TransactionWithSession(ctx, sess, tx1, tx2)
}

func tx1(_ context.Context, sess *xorm.Session) error {
	_, err := services.NewExampleService(sess).GetById(1)
	return err
}

func tx2(_ context.Context, sess *xorm.Session) error {
	_, err := services.NewExampleService(sess).GetById(2)
	return err
}
