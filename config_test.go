// author: wsfuyibing <websearch@163.com>
// date: 2022-06-06

package db

import (
    "encoding/json"
    "testing"
)

func TestConfiguration_Init(t *testing.T) {
    buf, _ := json.MarshalIndent(Config, "", "    ")
    t.Logf("config: %s.", buf)
}
