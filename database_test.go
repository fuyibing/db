// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-02-27

package db

import (
	"regexp"
	"testing"
)

func TestDatabase_GetHost(t *testing.T) {

	str := "user:passWord@123@tcp(host:3306)/mysql?charset=utf8"
	reg := regexp.MustCompile(`^([_a-zA-Z0-9-]+):(.+)@tcp\(([^)]+)\)/([^?]*)`)

	m := reg.FindStringSubmatch(str)

	for i := 1; i < len(m); i++ {
		t.Logf("arg %d: %s", i, m[i])
	}

	t.Logf("match: %v", m)

}
