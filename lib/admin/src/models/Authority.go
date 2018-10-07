package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

// jwt权限
type Authority struct {
	Id   int64
	User *User  `orm:"null;rel(fk)"`
	Name string `orm:"size(200)"` // jwt解开后的权限名称
}

func init() {
	orm.RegisterModel(new(Authority))
}

//更新用户
func UpdateUserAuth(auths []*Authority) (n int64, err error) {
	if len(auths) < 1 {
		return n, errors.New("集合为空")
	}
	o := orm.NewOrm()
	o.Begin()
	oneId := auths[0].User.Id
	_, err = o.Raw(`delete from authority where user_id=?`, oneId).Exec()
	if err != nil {
		o.Rollback()
		return n, err
	}
	tmp := Authority{}
	for _, v := range auths {
		tmp.Id = v.Id
		tmp.Name = v.Name
		tmp.User = &User{Id: v.User.Id}
		n, err = o.Insert(&tmp)
		if err != nil {
			return n, err
		}
	}
	o.Commit()
	return 1, nil
}
