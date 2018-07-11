package action

import (
	"auth/jwt/com"
	"io"
	"net/http"

	"github.com/bugfan/to"
)

const (
	LEN_USER = 5
	LEN_PWD  = 6
)

func Login(w http.ResponseWriter, r *http.Request) {

}

func Register(w http.ResponseWriter, r *http.Request) {
	obj := com.GetCtxMap(r)
	if err := com.HasValues(obj, "user", "password"); err != nil {
		io.WriteString(w, com.ToJsonString(com.Result{Status: 201, Msg: err.Error()}))
		return
	}
	user := to.String(obj["user"])
	password := to.String(obj["password"])
	if len(user) < LEN_USER {
		io.WriteString(w, com.ToJsonString(com.Result{Status: 202, Msg: "user less than " + to.String(LEN_USER)}))
		return
	}
	if len(password) < LEN_PWD {
		io.WriteString(w, com.ToJsonString(com.Result{Status: 202, Msg: "password less than " + to.String(LEN_PWD)}))
		return
	}
	o := com.NewOrm()
	o.Begin()
	q, _ := o.QueryTable("user").Filter("user", user).Count()
	if q > 0 {
		o.Commit()
		io.WriteString(w, com.ToJsonString(com.Result{Status: 203, Msg: "此用户已经注册"}))
		return
	}

}
