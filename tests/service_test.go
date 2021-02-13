// author: wsfuyibing <websearch@163.com>
// date: 2021-02-13

package tests

import (
	"testing"
	"time"

	"github.com/fuyibing/log"
	"xorm.io/xorm"

	"github.com/fuyibing/db"
)

func init() {
	log.Config.TimeFormat = "15:04:05.999999"
	log.Logger.SetAdapter(log.AdapterTerm)
	log.Logger.SetLevel(log.LevelDebug)
}

func TestService(t *testing.T) {

	ctx := log.NewContext()
	sess := db.SlaveContext(ctx)

	if err := db.TransactionWithSession(ctx, sess, dbHandle1); err != nil {
		log.Errorfc(ctx, "Transaction error: %v.", err)
		return
	}
	log.Infofc(ctx, "Transaction completed.")

	time.Sleep(time.Second)
}

func dbHandle1(ctx interface{}, sess *xorm.Session) error {
	if _, err := NewExampleService(sess).GetById(10); err != nil {
		return err
	}

	log.Infofc(ctx, "Not found")

	return nil
}

type Example struct {
	SubscriptionId int64
}

func (o *Example) TableName() string {
	return "mbs3_subscription"
}

type ExampleService struct {
	db.Service
}

func NewExampleService(s ...*xorm.Session) *ExampleService {
	o := &ExampleService{}
	o.Use(s...)
	return o
}

func (o *ExampleService) GetById(id int) (*Example, error) {
	model := &Example{}

	if _, err := o.Slave().Where("SubscriptionId = ?", id).Get(model); err != nil {
		return nil, err
	}
	if model.SubscriptionId > 0 {
		return model, nil
	}
	return nil, nil
}
