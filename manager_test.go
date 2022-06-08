// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

package db

import "testing"

func TestManagement_GetEngine(t *testing.T) {

    x := Manager.GetEngine("db")

    if x == nil {
        t.Error("engine not defined", x)
        return
    }

    buf, err := x.Slave().NewSession().QueryString("SELECT 1")
    if err != nil {
        t.Errorf("engine error: %v.", err)
        return
    }

    t.Logf("result: %s.", buf)
}
