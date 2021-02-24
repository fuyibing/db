// author: wsfuyibing <websearch@163.com>
// date: 2021-02-24

package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/fuyibing/log/v2"

	"github.com/fuyibing/db"
)

func TestPool(t *testing.T) {
	for i := 0; i < 100; i++ {
		loop()
		time.Sleep(time.Second * 5)
	}
}

func loop() {
	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go pool(wg, i)
	}
	wg.Wait()
}

func pool(wg *sync.WaitGroup, i int) {

	defer wg.Done()

	arr := &struct {
		Id int `xorm:"pk autoincr id"`
	}{}

	done, err := db.Master().SQL("SELECT * FROM test WHERE id = 100").Get(arr)
	if err != nil {
		log.Errorf("error get: %v.", err)
		return
	}

	log.Infof("ID: %v -> %v,", done, arr.Id)
}
