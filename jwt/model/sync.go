package model

import (
	"github.com/astaxie/beego/orm"
)

func Syncdb() {
	orm.RegisterModel(new(User))

	orm.RegisterModel()
}
