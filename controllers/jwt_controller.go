package controllers

import (
	"auth/lib/admin/src"
	com "auth/util"
	"auth/util/jwt"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/astaxie/beego/context"
	"github.com/bugfan/to"
)

// 建议规则
func Login(ctx *context.Context) {
	m := make(map[string]string)
	// err := json.Unmarshal(c.Ctx.Input.RequestBody, &m) beego和beego/context区别在于读取body
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(body, &m)
	log.Println("login", err, m)
	if err != nil {
		ctx.WriteString(com.ToJsonString(com.Result{
			Status: 401,
			Msg:    "请求参数不正确!",
		}))
		return
	}
	user, err := src.CheckLogin(m["username"], m["password"])
	if err != nil {
		ctx.WriteString(com.ToJsonString(com.Result{
			Status: 401,
			Msg:    "认证失败!", // 看了上野宣的 《图解http》,觉得应该不暴露详细错误了
		}))
		return
	}
	authMap := make(map[string]interface{})
	authMap["Username"] = user.Username
	authMap["Id"] = user.Id
	authMap["Authority"] = user.Authority
	authMap["Email"] = user.Email
	authMap["Nickname"] = user.Nickname
	authBs, err := json.Marshal(authMap)
	jwtStr, err := jwt.GetJWT(string(authBs))
	if err != nil {
		ctx.WriteString(com.ToJsonString(com.Result{
			Status: 401,
			Msg:    to.String(err),
		}))
		return
	}
	ctx.WriteString(com.ToJsonString(com.Result{
		Status: 200,
		Data:   to.String(jwtStr),
		Msg:    "成功",
	}))
	return
}

func ChangePassword(ctx *context.Context) {
	ctx.WriteString(com.ToJsonString(com.Result{
		Status: 200,
		Msg:    "成功",
	}))
}
func Logout(ctx *context.Context) {
	ctx.WriteString(com.ToJsonString(com.Result{
		Status: 200,
		Msg:    "成功",
	}))
}
