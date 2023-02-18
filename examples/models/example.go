// author: wsfuyibing <websearch@163.com>
// date: 2023-02-18

package models

type Example struct {
	Id int `xorm:"id pk autoincr"`
}

func (o *Example) TableName() string {
	return "task"
}
