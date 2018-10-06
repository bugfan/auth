package rbac

import (
	"auth/lib/admin/src/models"
	m "auth/lib/admin/src/models"

	"github.com/astaxie/beego/orm"
)

type UserController struct {
	CommonController
}

func (this *UserController) Index() {
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("rows")
	sort := this.GetString("sort")
	order := this.GetString("order")
	if len(order) > 0 {
		if order == "desc" {
			sort = "-" + sort
		}
	} else {
		sort = "Id"
	}
	users, count := m.Getuserlist(page, page_size, sort)
	users2 := []orm.Params{}
	for _, v := range users {
		v["Authority"] = []models.Authority{models.Authority{Name: "test"}, models.Authority{Name: "admin"}, models.Authority{Name: "hh"}}
		users2 = append(users2, v)
	}
	if this.IsAjax() {
		this.Data["json"] = &map[string]interface{}{"total": count, "rows": &users2}
		this.ServeJSON()
		return
	} else {
		tree := this.GetTree()
		this.Data["tree"] = &tree
		this.Data["users"] = &users2
		if this.GetTemplatetype() != "easyui" {
			this.Layout = this.GetTemplatetype() + "/public/layout.tpl"
		}
		this.TplName = this.GetTemplatetype() + "/rbac/user.tpl"
	}

}

func (this *UserController) AddUser() {
	u := m.User{}
	if err := this.ParseForm(&u); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	id, err := m.AddUser(&u)
	if err == nil && id > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *UserController) UpdateUser() {
	u := m.User{}
	if err := this.ParseForm(&u); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	id, err := m.UpdateUser(&u)
	if err == nil && id > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *UserController) DelUser() {
	Id, _ := this.GetInt64("Id")
	status, err := m.DelUserById(Id)
	if err == nil && status > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}
