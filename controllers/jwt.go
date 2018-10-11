package controllers

import (
	"auth/lib/admin/src"
	"encoding/json"

	"github.com/astaxie/beego"
)

type JWTController struct {
	beego.Controller
}

func (c *JWTController) Post() {
	defer c.ServeJSON(true)
	m := make(map[string]string)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	if err != nil {
		c.Data["Status"] = 401
		c.Data["Msg"] = "请求参数不正确!"
		return
	}
	user, err := src.CheckLogin(m["username"], m["password"])
	if err != nil {
		c.Data["Status"] = 401
		c.Data["Msg"] = err
		return
	}

}
