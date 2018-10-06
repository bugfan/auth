package models

import "github.com/astaxie/beego/orm"

// jwt权限
type Authority struct {
	Id   int64
	User *User  `orm:"null;rel(fk)"`
	Name string `orm:"size(200)"` // jwt解开后的权限名称
}

func init() {
	orm.RegisterModel(new(Authority))
}
